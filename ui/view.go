package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// cellW is the single grid unit shared by pyramid columns AND multiplier bins.
const cellW = 7

// bWidth = 17 slots × cellW. Changing cellW keeps everything aligned.
const bWidth = 17 * cellW // 119

const sWidth = 28

// bl renders content at exactly bWidth with pure-black background.
// Every board line goes through this — eliminates gray bleed.
func bl(content string) string {
	return lipgloss.NewStyle().
		Width(bWidth).
		Background(ColorPureBlack).
		Render(content)
}

// pinCell centers a rendered string of known visual length inside a cellW slot.
func pinCell(rendered string, visualLen int) string {
	left := (cellW - visualLen) / 2
	right := cellW - visualLen - left
	return strings.Repeat(" ", left) + rendered + strings.Repeat(" ", right)
}

func (m Model) View() string {
	if m.state == StateVerify {
		return m.renderVerify()
	}

	var board, side strings.Builder

	line := func(s string) { board.WriteString(bl(s) + "\n") }
	blank := func() { board.WriteString(bl("") + "\n") }
	div := func() { line(MutedStyle.Render(strings.Repeat("─", bWidth))) }

	// ── TITLE ────────────────────────────────────────────────────────────────
	stateTag := GoldStyle.Render(" ◆ READY ")
	if m.state == StateDropping {
		stateTag = PinkStyle.Render(" ◆ DROPPING ")
	}
	line(lipgloss.NewStyle().Width(bWidth).Align(lipgloss.Center).
		Render(TitleStyle.Render(" CIPHER PLINKO PRO ") + "  " + stateTag))
	blank()

	// ── RANK / XP ────────────────────────────────────────────────────────────
	xpMax := m.level * 200
	xpPct := 0
	if xpMax > 0 {
		xpPct = (m.xp * 40) / xpMax
		if xpPct > 40 {
			xpPct = 40
		}
	}
	xpBar := CyanStyle.Render(strings.Repeat("█", xpPct)) +
		MutedStyle.Render(strings.Repeat("░", 40-xpPct))
	line(fmt.Sprintf("  RANK: %s   LVL: %s   [ %s ]",
		GoldStyle.Render(fmt.Sprintf("%-12s", m.rank)),
		GoldStyle.Render(fmt.Sprintf("%d", m.level)),
		xpBar,
	))
	blank()

	// ── STATS BAR ────────────────────────────────────────────────────────────
	autoStr := MutedStyle.Render("OFF")
	if m.autoBot {
		autoStr = GoldStyle.Render("ON ")
	}
	line(fmt.Sprintf("  BAL: %s     BET: %s × %s     RISK: %s     AUTO: %s",
		GoldStyle.Render(fmt.Sprintf("$%.2f", m.balance)),
		PinkStyle.Render(fmt.Sprintf("$%.2f", m.stake)),
		CyanStyle.Render(fmt.Sprintf("%d", m.ballCount)),
		CyanStyle.Render(m.riskLevel),
		autoStr,
	))
	div()
	blank()

	// ── PYRAMID ──────────────────────────────────────────────────────────────
	for r := 0; r <= m.rows; r++ {
		var row strings.Builder
		row.WriteString(strings.Repeat(" ", (m.rows-r)*cellW/2))
		for c := 0; c <= r; c++ {
			char, style := "·", MutedStyle
			for _, b := range m.balls {
				if !b.done && b.waitTicks == 0 && b.currentRow == r && b.currentCol == c {
					char, style = "●", PinkStyle
				}
			}
			row.WriteString(pinCell(style.Render(char), 1))
		}
		line(row.String())
	}

	// ── MULTIPLIER BINS ───────────────────────────────────────────────────────
	blank()
	var bins strings.Builder
	for _, v := range m.riskMults[m.riskLevel] {
		label := fmt.Sprintf("[%g]", v)
		var s lipgloss.Style
		switch {
		case v >= 10:
			s = GoldStyle
		case v >= 1:
			s = lipgloss.NewStyle().Foreground(ColorNeonGreen).Bold(true)
		default:
			s = MutedStyle
		}
		bins.WriteString(pinCell(s.Render(label), len(label)))
	}
	line(bins.String())

	// ── SIDE PANEL ────────────────────────────────────────────────────────────
	// HISTORY
	side.WriteString(HeaderStyle.Render("  HISTORY") + "\n\n")
	const historySlots = 8
	for i := 0; i < historySlots; i++ {
		if i < len(m.history) {
			h := m.history[i]
			isLoss := strings.HasPrefix(h, "-")
			sym := GoldStyle.Render("▲")
			es := WinStyle
			if isLoss {
				sym = LossStyle.Render("▼")
				es = LossStyle
			}
			side.WriteString(fmt.Sprintf("  %s %s\n", sym, es.Render(h)))
		} else {
			side.WriteString("\n")
		}
	}

	// CONTROLS
	side.WriteString("\n" + MutedStyle.Render(strings.Repeat("─", sWidth-2)) + "\n\n")
	side.WriteString(HeaderStyle.Render("  CONTROLS") + "\n\n")
	for _, c := range [][2]string{
		{"ENTER", "Drop Ball"},
		{"↑ / ↓", "Bet Size"},
		{"← / →", "Ball Count"},
		{"R", "Risk Level"},
		{"B", "Auto-Bot"},
		{"V", "Verify"},
		{"Q", "Quit"},
	} {
		side.WriteString(fmt.Sprintf("  %s  %s\n",
			CyanStyle.Render(fmt.Sprintf("%-7s", c[0])),
			MutedStyle.Render(c[1])))
	}

	// SESSION STATS — fills the height gap between board and side panel
	side.WriteString("\n" + MutedStyle.Render(strings.Repeat("─", sWidth-2)) + "\n\n")
	side.WriteString(HeaderStyle.Render("  SESSION") + "\n\n")

	pl := m.balance - m.startBalance
	plStr := WinStyle.Render(fmt.Sprintf("+$%.2f", pl))
	if pl < 0 {
		plStr = LossStyle.Render(fmt.Sprintf("-$%.2f", math.Abs(pl)))
	}
	side.WriteString(fmt.Sprintf("  %s  %s\n",
		MutedStyle.Render("Rounds :"),
		CyanStyle.Render(fmt.Sprintf("%d", m.sessionRounds))))
	side.WriteString(fmt.Sprintf("  %s  %s\n",
		MutedStyle.Render("P&L    :"),
		plStr))
	side.WriteString(fmt.Sprintf("  %s  %s\n",
		MutedStyle.Render("Nonce  :"),
		MutedStyle.Render(fmt.Sprintf("%d", m.engine.Nonce))))
	if m.lastResult != "" {
		side.WriteString("\n")
		isLoss := strings.HasPrefix(m.lastResult, "-")
		if isLoss {
			side.WriteString(fmt.Sprintf("  %s\n", LossStyle.Render(m.lastResult)))
		} else {
			side.WriteString(fmt.Sprintf("  %s\n", WinStyle.Render(m.lastResult)))
		}
	}

	layout := lipgloss.JoinHorizontal(lipgloss.Top,
		board.String(),
		lipgloss.NewStyle().Width(sWidth).Background(ColorPureBlack).Render(side.String()),
	)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		MainFrame.Render(layout))
}

func (m Model) renderVerify() string {
	var sb strings.Builder
	sb.WriteString(TitleStyle.Render(" CIPHER VERIFY — PROVABLY FAIR ") + "\n\n")
	sb.WriteString(fmt.Sprintf("  %s  %s\n",
		CyanStyle.Render("Server Seed :"), MutedStyle.Render(m.engine.ServerSeed)))
	sb.WriteString(fmt.Sprintf("  %s  %s\n",
		CyanStyle.Render("Client Seed :"), MutedStyle.Render(m.engine.ClientSeed)))
	sb.WriteString(fmt.Sprintf("  %s  %s\n\n",
		CyanStyle.Render("Nonce       :"), GoldStyle.Render(fmt.Sprintf("%d", m.engine.Nonce))))

	sb.WriteString(MutedStyle.Render(strings.Repeat("─", bWidth)) + "\n\n")
	sb.WriteString(HeaderStyle.Render("AUDIT LOG — LAST 5 ROUNDS") + "\n\n")

	logs := m.auditLogs
	if len(logs) > 5 {
		logs = logs[len(logs)-5:]
	}
	if len(logs) == 0 {
		sb.WriteString(MutedStyle.Render("  No rounds played yet.\n"))
	}
	for i, l := range logs {
		winStr := WinStyle.Render(fmt.Sprintf("$%.2f", l.Win))
		if l.Win == 0 {
			winStr = LossStyle.Render("$0.00")
		}
		sb.WriteString(fmt.Sprintf("  %s   WIN: %s\n",
			GoldStyle.Render(fmt.Sprintf("#%d", i+1)), winStr))
		sb.WriteString(MutedStyle.Render(fmt.Sprintf("         Hash  : %s...\n", l.Hash[:24])))
		sb.WriteString(MutedStyle.Render(fmt.Sprintf("         Nonce : %d\n\n", l.Nonce)))
	}
	sb.WriteString("\n" + MutedStyle.Render("  Press [V] to return to game."))

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		MainFrame.Render(lipgloss.NewStyle().Width(bWidth).Background(ColorPureBlack).
			Render(sb.String())))
}

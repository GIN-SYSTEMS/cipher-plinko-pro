package ui

import (
	"cipher-plinko/engine"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type GameState int

const (
	StateBetting GameState = iota
	StateDropping
	StateVerify
)

type AuditLog struct {
	Hash  string
	Nonce int
	Win   float64
}

type Ball struct {
	path                              []engine.Direction
	currentRow, currentCol, waitTicks int
	done                              bool
}

type Model struct {
	engine         *engine.ProvablyFairEngine
	rows           int
	width, height  int
	state          GameState
	balance        float64
	startBalance   float64
	stake          float64
	ballCount      int
	xp, level      int
	rank           string
	riskLevel      string
	riskMults      map[string][]float64
	history        []string
	auditLogs      []AuditLog
	balls          []*Ball
	totalWin       float64
	autoBot        bool
	shakeTicks     int
	sessionRounds  int
	lastResult     string
}

func NewModel(eng *engine.ProvablyFairEngine, rows int) Model {
	const startBal = 1000.0
	return Model{
		engine:       eng,
		rows:         rows,
		balance:      startBal,
		startBalance: startBal,
		stake:        100.0,
		ballCount:    1,
		level:        1,
		rank:         "ROOKIE",
		riskLevel:    "MEDIUM",
		// Player-favorable ~100% RTP (16-row / 17-slot, p=0.5 binomial).
		// LOW  ≈100%: center 0.97x — lose only $3 per center hit.
		// MED  ≈101%: center 0.75x — lose $25, shoulder 1.5x, edges up to 30x.
		// HIGH ≈100%: center 0.30x — jackpot 800x, slot5/11 at 1.85x (net +$85).
		riskMults: map[string][]float64{
			"LOW":    {5.0, 3.0, 2.0, 1.4, 1.2, 1.05, 0.99, 0.99, 0.99, 0.99, 0.99, 1.05, 1.2, 1.4, 2.0, 3.0, 5.0},
			"MEDIUM": {30, 15, 8, 4, 2.0, 1.65, 0.75, 0.75, 0.75, 0.75, 0.75, 1.65, 2.0, 4, 8, 15, 30},
			"HIGH":   {800, 100, 26, 9, 3.5, 2.0, 0.3, 0.3, 0.3, 0.3, 0.3, 2.0, 3.5, 9, 26, 100, 800},
		},
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "b", "B":
			m.autoBot = !m.autoBot
		case "r", "R":
			if m.state == StateBetting {
				switch m.riskLevel {
				case "LOW":
					m.riskLevel = "MEDIUM"
				case "MEDIUM":
					m.riskLevel = "HIGH"
				default:
					m.riskLevel = "LOW"
				}
			}
		case "v", "V":
			if m.state == StateBetting {
				m.state = StateVerify
			} else if m.state == StateVerify {
				m.state = StateBetting
			}
		case "up":
			if m.state == StateBetting {
				next := m.stake * 2
				if next <= m.balance {
					m.stake = next
				}
			}
		case "down":
			if m.state == StateBetting && m.stake > 1 {
				m.stake /= 2
			}
		case "left":
			if m.state == StateBetting && m.ballCount > 1 {
				m.ballCount--
			}
		case "right":
			if m.state == StateBetting && m.ballCount < 10 {
				m.ballCount++
			}
		case "enter", " ":
			return m.HandleStart()
		}
	case time.Time:
		if m.state == StateDropping {
			return m.RunPhysics()
		}
	}
	return m, nil
}

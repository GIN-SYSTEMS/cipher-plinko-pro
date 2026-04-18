package ui

import (
	"fmt"
	"time"

	"cipher-plinko/engine"

	tea "github.com/charmbracelet/bubbletea"
)

func tickEvery() tea.Cmd {
	return tea.Tick(30*time.Millisecond, func(t time.Time) tea.Msg { return t })
}

func (m Model) HandleStart() (tea.Model, tea.Cmd) {
	if m.state != StateBetting {
		return m, nil
	}
	cost := m.stake * float64(m.ballCount)
	if m.balance < cost {
		return m, nil
	}

	m.balance -= cost
	m.balls = nil
	m.totalWin = 0

	for i := 0; i < m.ballCount; i++ {
		m.engine.Nonce++
		m.balls = append(m.balls, &Ball{
			path:      m.engine.CalculatePath(m.rows),
			waitTicks: i * 3,
		})
	}

	m.state = StateDropping
	return m, tickEvery()
}

func (m Model) RunPhysics() (tea.Model, tea.Cmd) {
	allDone := true

	for _, b := range m.balls {
		if b.done {
			continue
		}
		if b.waitTicks > 0 {
			b.waitTicks--
			allDone = false
			continue
		}
		if b.currentRow < m.rows {
			if b.path[b.currentRow] == engine.Right {
				b.currentCol++
			}
			b.currentRow++
			b.waitTicks = 1
			allDone = false
		} else {
			slot := b.currentCol
			mults := m.riskMults[m.riskLevel]
			if slot >= len(mults) {
				slot = len(mults) - 1
			}
			m.totalWin += m.stake * mults[slot]
			b.done = true
		}
	}

	if allDone {
		return m.settle()
	}
	return m, tickEvery()
}

func (m Model) settle() (tea.Model, tea.Cmd) {
	m.balance += m.totalWin
	cost := m.stake * float64(m.ballCount)
	net := m.totalWin - cost

	if m.totalWin > 0 {
		m.xp += int(m.totalWin)
	}
	for m.xp >= m.level*200 {
		m.xp -= m.level * 200
		m.level++
	}

	switch {
	case m.level >= 50:
		m.rank = "CIPHER ELITE"
	case m.level >= 20:
		m.rank = "WHALE"
	case m.level >= 5:
		m.rank = "PRO"
	default:
		m.rank = "ROOKIE"
	}

	m.auditLogs = append(m.auditLogs, AuditLog{
		Hash:  m.engine.GenerateHash(),
		Nonce: m.engine.Nonce,
		Win:   m.totalWin,
	})

	m.sessionRounds++

	mult := 0.0
	if m.ballCount > 0 && m.stake > 0 {
		mult = m.totalWin / (m.stake * float64(m.ballCount))
	}
	var entry string
	if net >= 0 {
		entry = fmt.Sprintf("+$%.2f  ×%.4g", net, mult)
		m.lastResult = entry
	} else {
		entry = fmt.Sprintf("-$%.2f  ×%.4g", -net, mult)
		m.lastResult = entry
	}
	m.history = append([]string{entry}, m.history...)
	if len(m.history) > 8 {
		m.history = m.history[:8]
	}

	m.state = StateBetting
	if m.autoBot && m.balance >= m.stake*float64(m.ballCount) {
		return m.HandleStart()
	}
	return m, nil
}

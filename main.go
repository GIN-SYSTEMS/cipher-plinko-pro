package main

import (
	"fmt"
	"os"

	"cipher-plinko/engine"
	"cipher-plinko/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	serverSeed := "super_secret_server_seed_2026"
	clientSeed := "player_custom_seed"
	nonce := 1
	rows := 16

	gameEngine := engine.NewEngine(serverSeed, clientSeed, nonce)
	m := ui.NewModel(gameEngine, rows)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("failed to start game: %v", err)
		os.Exit(1)
	}
}

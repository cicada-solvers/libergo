package main

import (
	"fmt"
	"math/big"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	counter  *big.Int
	quitting bool
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		return m, tickCmd()
	case counterMsg:
		m.counter = msg.counter
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "q" {
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}
	return fmt.Sprintf("Current counter: %s\nPress 'q' to quit.\n", m.counter.String())
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type counterMsg struct {
	counter *big.Int
}

func updateCounterCmd(counter *big.Int) tea.Cmd {
	return func() tea.Msg {
		return counterMsg{counter: counter}
	}
}

package ssh

import (
	"github.com/ahmedyahia/deathmatch/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

func TeaHandler() bubbletea.Handler {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		return ui.NewModel(), []tea.ProgramOption{tea.WithAltScreen()}
	}
}

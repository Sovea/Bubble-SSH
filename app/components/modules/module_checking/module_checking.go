package module_checking

import (
	"fmt"
	// "os"
	"bubbles_ssh/app/service/service_network"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error
type net_checkMsg struct {
	status int
	err    error
	msg    string
}
type model_loading struct {
	spinner  spinner.Model
	status   int
	quitting bool
	err      error
	msg      string
}

// type model_loading_msg struct {
// 	model  model_loading
// 	msg string
// }
func InitialModel(default_msg string, spinner_style spinner.Spinner) model_loading {
	s := spinner.NewModel()
	s.Spinner = spinner_style
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model_loading{spinner: s, msg: default_msg, status: 0}
}

func (m model_loading) Init() tea.Cmd {
	return tea.Batch(Network_Check(), spinner.Tick)
}

func (m model_loading) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	case net_checkMsg:
		if msg.status == 0 {
			return m, tick()
		} else if msg.status == 1 {
			m.status = 1
			m.msg = msg.msg
			return m, tea.Quit
		} else {
			m.status = 2
			m.msg = msg.msg
			return m, tea.Quit
		}

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

}

func (m model_loading) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	// fmt.Println(&m)
	str := fmt.Sprintf("\n\n	%s   %s \n\n", m.spinner.View(), m.msg)
	if m.quitting {
		return str + "\n"
	}
	return str
}

func tick() tea.Cmd {
	return spinner.Tick
}

func Network_Check() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		network_status := service_network.NetWorkStatus()
		if network_status == true {
			// m.msg = "Network is avaliable."
			// m.status = 1
			return net_checkMsg{status: 1, msg: "Network is avaliable."}
		} else {
			// m.msg = "Network is not avaliable. Please fix it before start bubble-ssh."
			// m.status = 2
			return net_checkMsg{status: 2, msg: "Network is not avaliable. Please fix it before start bubble-ssh."}
		}
	})

}


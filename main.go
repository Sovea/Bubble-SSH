package main

import (
	"bubbles_ssh/app/components/modules/module_checking"
	"bubbles_ssh/app/components/modules/module_chose_env"
	"bubbles_ssh/app/components/modules/module_page"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// "bubbles_ssh/app/service/service_network"
)

func Intro() {
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingTop(2).
		PaddingLeft(4).
		Width(22)
	fmt.Println(style.Render("Hello, kitty."))
}
func main() {
	// Intro()
	module_page.Page()
	instance_loading := module_checking.InitialModel("Network Check Task Dispatching ...", spinner.Points)
	program_net_check := tea.NewProgram(instance_loading)
	if err := program_net_check.Start(); err != nil {
		fmt.Printf("Could not start program...\n\n Please check user permissions and network status, etc...")
		os.Exit(1)
	}
	instance_chose_env := module_chose_env.InitialModel()
	program_chose_env := tea.NewProgram(instance_chose_env)
	if err := program_chose_env.Start(); err != nil {
		fmt.Printf("Could not start ssh handle program... \n\n Please check user permissions, software source and network status, etc...")
		os.Exit(1)
	}
}

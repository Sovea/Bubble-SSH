package module_chose_env

import (
	"bubbles_ssh/app/service/service_cmd"
	"bubbles_ssh/app/service/service_port"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	// "github.com/fogleman/ease"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
)

// General stuff for styling the view
var (
	term          = termenv.ColorProfile()
	keyword       = makeFgStyle("211")
	subtle        = makeFgStyle("241")
	progressEmpty = subtle(progressEmptyChar)
	dot           = colorFg(" • ", "236")

	// Gradient colors we'll use for the progress bar
	ramp             = makeRamp("#B14FFF", "#00FFA3", progressBarWidth)
	env_List         = []string{"Ubuntu", "CentOS", "The Others"}
	env_install_List = [][]string{{"ssh", "sshd_config", "ssh_service"}, {"ssh", "sshd_config", "ssh_service"}, {}}
	chat_message     = []string{"Cool, we need the ", "OKey, we should install ", "Oh sorry, we do not support other versions of Linux at the moment."}
)

type tickMsg struct{}
type frameMsg struct{}
type setupMsg struct {
	name       string
	status     int
	err        error
	suggestion string
}
type exitMsg struct {
	msg        string
	suggestion string
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

func (m chose_env_model) exit_handler() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return exitMsg{}
	})
}

type chose_env_model struct {
	spinner   spinner.Model
	Choice    int
	Chosen    bool
	Ticks     int
	setupStep int
	Quitting  bool
	exitmsg   exitMsg
}

// type model_loading struct {
// 	spinner  spinner.Model
// 	status   int
// 	quitting bool
// 	err      error
// 	msg      string
// }

func (m chose_env_model) ssh_handler() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		port_status := service_port.Raw_connect("127.0.0.1", []string{"22"})
		ssh_handler_status := false
		var err error = nil
		if port_status == false {
			if m.Choice == 0 {
				ssh_handler_status, err = service_cmd.CMD_Exec("apt", "install", "openssh-server", "-y")
			}
			if ssh_handler_status == true {
				return setupMsg{name: "ssh", status: 1, err: nil}
			} else {
				return setupMsg{name: "ssh", status: 2, err: err, suggestion: "openssh-server安装失败，请检查软件源和包依赖，建议使用" + keyword("lsb_release") + "等命令查看发行版本并更换正确的软件源。\n\n The openssh-server installation fails. Please check the " + keyword("software source") + " and " + keyword("package dependencies") + ". It is recommended to use commands such as " + keyword("lsb_release") + " to check the release version and replace the correct software source."}
			}
		}
		return setupMsg{name: "ssh", status: 2, err: fmt.Errorf("The ssh service exists."), suggestion: ""}
	})
}
func (m chose_env_model) ssh_service_handler() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		ssh_service_status := false
		var err error = nil
		if m.Choice == 0 {
			ssh_service_status, err = service_cmd.CMD_Exec("service", "ssh", "restart")
		}
		if ssh_service_status == true {
			return setupMsg{name: "ssh_service", status: 1, err: nil}
		} else {
			ssh_service_status, err = service_cmd.CMD_Exec("/etc/init.d/ssh", "restart")
			if ssh_service_status == true {
				return setupMsg{name: "ssh_service", status: 1, err: nil}
			} else {
				return setupMsg{name: "ssh_service", status: 2, err: err, suggestion: "ssh服务重启时失败，请检查一下.\n\nThe ssh service failed when restarting, please check it."}
			}
		}
	})
}
func WriteToFile(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		// fmt.Println("file create failed. err: " + err.Error())
	} else {
		// offset
		//os.Truncate(filename, 0) //clear
		n, _ := f.Seek(0, os.SEEK_END)
		_, err = f.WriteAt([]byte(content), n)
		// fmt.Println("write succeed!")
		defer f.Close()
	}
	return err
}
func defer_close(defer_time time.Duration) tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		time.Sleep(defer_time)
		return tea.Quit()
	})
}
func (m chose_env_model) sshd_config_handler() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		sshd_handler_status := false
		var err error = nil
		if m.Choice == 0 {
			sshd_handler_status, err = service_cmd.CMD_Exec("mv", "/etc/ssh/sshd_config", "/etc/ssh/sshd_config.bubble_back")
			if sshd_handler_status == true {
				err = WriteToFile("/etc/ssh/sshd_config", `# This is the sshd server system-wide configuration file.  

# The strategy used for options in the default sshd_config shipped with
# OpenSSH is to specify options with their default value where
# possible, but leave them commented.  Uncommented options override the
# default value.

#Port 22
#AddressFamily any
AddressFamily inet
#ListenAddress 0.0.0.0
#ListenAddress ::

#HostKey /etc/ssh/ssh_host_rsa_key
#HostKey /etc/ssh/ssh_host_ecdsa_key
#HostKey /etc/ssh/ssh_host_ed25519_key

# Ciphers and keying
#RekeyLimit default none

# Logging
#SyslogFacility AUTH
#LogLevel INFO

# Authentication:

#LoginGraceTime 2m
#PermitRootLogin yes
#StrictModes yes
#MaxAuthTries 6
#MaxSessions 10

#PubkeyAuthentication yes

# Expect .ssh/authorized_keys2 to be disregarded by default in future.
#AuthorizedKeysFile     .ssh/authorized_keys .ssh/authorized_keys2

#AuthorizedPrincipalsFile none

#AuthorizedKeysCommand none
#AuthorizedKeysCommandUser nobody

# For this to work you will also need host keys in /etc/ssh/ssh_known_hosts
#HostbasedAuthentication no
# Change to yes if you don't trust ~/.ssh/known_hosts for
# HostbasedAuthentication
#IgnoreUserKnownHosts no
# Don't read the user's ~/.rhosts and ~/.shosts files
#IgnoreRhosts yes

# To disable tunneled clear text passwords, change to no here!
PasswordAuthentication yes
#PermitEmptyPasswords no

# Change to yes to enable challenge-response passwords (beware issues with
# some PAM modules and threads)
ChallengeResponseAuthentication no

# Kerberos options
#KerberosAuthentication no
#KerberosOrLocalPasswd yes
#KerberosTicketCleanup yes
#KerberosGetAFSToken no

# GSSAPI options
#GSSAPIAuthentication no
#GSSAPICleanupCredentials yes
#GSSAPIStrictAcceptorCheck yes
#GSSAPIKeyExchange no

# Set this to 'yes' to enable PAM authentication, account processing,
# and session processing. If this is enabled, PAM authentication will
# be allowed through the ChallengeResponseAuthentication and
# PasswordAuthentication.  Depending on your PAM configuration,
# PAM authentication via ChallengeResponseAuthentication may bypass
# the setting of "PermitRootLogin without-password".
# If you just want the PAM account and session checks to run without
# PAM authentication, then enable this but set PasswordAuthentication
# and ChallengeResponseAuthentication to 'no'.
UsePAM yes

#AllowAgentForwarding yes
#AllowTcpForwarding yes
#GatewayPorts no
X11Forwarding yes
#X11DisplayOffset 10
#X11UseLocalhost yes
#PermitTTY yes
PrintMotd no
#PrintLastLog yes
#TCPKeepAlive yes
#UseLogin no
#PermitUserEnvironment no
#Compression delayed
#ClientAliveInterval 0
#ClientAliveCountMax 3
#UseDNS no
#PidFile /var/run/sshd.pid
#MaxStartups 10:30:100
#PermitTunnel no
#ChrootDirectory none
#VersionAddendum none

# no default banner path
#Banner none

# Allow client to pass locale environment variables
AcceptEnv LANG LC_*

# override default of no subsystems
Subsystem sftp  /usr/lib/openssh/sftp-server

# Example of overriding settings on a per-user basis
#Match User anoncvs
#       X11Forwarding no
#       AllowTcpForwarding no
#       PermitTTY no
#       ForceCommand cvs server
`)
				if err != nil {
					return setupMsg{name: "sshd_config", status: 2, err: err}
				}
				return setupMsg{name: "sshd_config", status: 1, err: nil}
			} else {
				return setupMsg{name: "sshd_config", status: 2, err: err}
			}
		}
		return setupMsg{name: "sshd_config", status: 2, err: err, suggestion: "Failed to write configuration file, please check your sshd_config file."}
	})
}

func InitialModel() chose_env_model {
	s := spinner.NewModel()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return chose_env_model{spinner: s, Choice: 0, Chosen: false, Ticks: 15, setupStep: 0, Quitting: false}
}
func (m chose_env_model) Init() tea.Cmd {
	return tick()
}

// Main update function.
func (m chose_env_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.Chosen {
		return updateChoices(msg, m)
	}

	return updateChosen(msg, m)
}

// The main view, which just calls the appropriate sub-view
func (m chose_env_model) View() string {
	var s string
	// if m.Quitting {
	// 	return "\n  See you later!\n\n"
	// }
	if !m.Chosen {
		s = choicesView(m)
	} else {
		if m.setupStep == 0 {
			s = chosenView(m)
		} else {
			// time.Sleep(3)
			s = setupView(m)
		}
	}
	return indent.String("\n"+s+"\n\n", 2)
}

// Update loop for the first view where you're choosing the env.
func updateChoices(msg tea.Msg, m chose_env_model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice += 1
			if m.Choice >= len(env_List) {
				m.Choice = len(env_List) - 1
			}
		case "k", "up":
			m.Choice -= 1
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			m.Ticks = 0
			if m.Choice < (len(env_List) - 1) {
				return m, frame()
			} else {
				return m, tea.Quit
			}

		}

	case tickMsg:
		if m.Ticks == 0 {
			m.Quitting = true
			return m, tea.Quit
		}
		m.Ticks -= 1
		return m, tick()
	}

	return m, nil
}

// Update loop for the second view after a choice has been made
func updateChosen(msg tea.Msg, m chose_env_model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case frameMsg:
		if m.setupStep == 0 {
			m.setupStep += 1
			return m, tea.Batch(spinner.Tick, m.ssh_handler())
		}

	case setupMsg:
		if msg.status == 1 {
			if m.setupStep == 1 {
				m.setupStep += 1
				return m, tea.Batch(spinner.Tick, m.sshd_config_handler())
			} else if m.setupStep == 2 {
				m.setupStep += 1
				return m, tea.Batch(spinner.Tick, m.ssh_service_handler())
			} else if m.setupStep == 3 {
				m.setupStep += 1
				return m, tea.Batch(spinner.Tick, defer_close(3))
			}
		} else {
			// fmt.Println(msg.err)
			m.exitmsg = exitMsg{msg: msg.err.Error(), suggestion: msg.suggestion}
			return m, tea.Batch(m.exit_handler(), defer_close(3))
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case exitMsg:
		return m, tea.Batch(m.exit_handler(), defer_close(3))

	default:
		return m, spinner.Tick
	}
	return m, spinner.Tick
}

// The choices view, when you're choosing a task
func choicesView(m chose_env_model) string {
	c := m.Choice

	tpl := "Please choose your Linux distribution:\n\n"
	tpl += "%s\n\n"
	tpl += "Program will quits in %s seconds\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")
	choice_list := ""
	for k, v := range env_List {
		choice_list += checkbox(v, c == k) + "\n"
	}
	// choices := fmt.Sprintf(
	// 	"%s\n%s\n%s\n",
	// 	checkbox(env_List[0], c == 0),
	// 	checkbox(env_List[1], c == 1),
	// 	checkbox("The Others", c == 2),
	// )

	return fmt.Sprintf(tpl, choice_list, colorFg(strconv.Itoa(m.Ticks), "79"))
}

// The quick chosen view, after a task has been chosen
func chosenView(m chose_env_model) string {
	var msg string
	msg = env_List[m.Choice] + "?\n\n" + chat_message[m.Choice]
	for _, v := range env_install_List[m.Choice] {
		msg += keyword(v)
		msg += " "
	}
	msg += "..."
	label := ""

	return msg + "\n\n" + label + "\n"
}

// The install view, after a task has been chosen ang gonna install something
func setupView(m chose_env_model) string {
	// fmt.Println(&m)
	var str string
	if m.exitmsg.msg != "" {
		str = m.exitmsg.msg + "\n\n" + m.exitmsg.suggestion + "\n\n"
	} else {
		if m.setupStep < len(env_install_List[m.Choice]) {
			str = fmt.Sprintf("\n\n	%s Check and handle the %s ...\n\n", m.spinner.View(), env_install_List[m.Choice][m.setupStep-1])
		} else {
			str = fmt.Sprintf("\n\n The SSH service is avaliable, please check it out ...\n\n")
		}
	}

	return str
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

func progressbar(width int, percent float64) string {
	w := float64(progressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += termenv.String(progressFullChar).Foreground(term.Color(ramp[i])).String()
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

// Utils

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

// Color a string's foreground and background with the given value.
func makeFgBgStyle(fg, bg string) func(string) string {
	return termenv.Style{}.
		Foreground(term.Color(fg)).
		Background(term.Color(bg)).
		Styled
}

// Generate a blend of colors.
func makeRamp(colorA, colorB string, steps float64) (s []string) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, colorToHex(c))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format compatible with termenv.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}

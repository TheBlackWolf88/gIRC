package main

import (
	"bufio"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"net"
	"os"
)

type model struct {
	user  string
	chat  string
	conn  net.Conn
	input textinput.Model
}

func main() {
	_, err := tea.NewProgram(initialModel()).Run()
	if err != nil {
		os.Exit(1)
	}
}

func initialModel() model {
	conn, _ := net.Dial("tcp", "localhost:8080")
	ti := textinput.New()
	ti.Placeholder = "Say something!"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return model{
		user:  os.Getenv("$USER"),
		chat:  "",
		conn:  conn,
		input: ti,
	}
}

type chatMsg []byte

func checkServer(m model) tea.Cmd {
	return m.readServer
}

func (m model) readServer() tea.Msg {
	//listen for reply
	for {
		message, _ := bufio.NewReader(m.conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
		return chatMsg(message)
	}
}

func (m model) writeServer() tea.Msg {
	//listen for reply
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		//send to socket
		fmt.Fprint(m.conn, m.user+": "+text)
	}
}

func (m model) Init() tea.Cmd {
	return checkServer(m)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case chatMsg:
		m.chat = m.chat + string(msg)
		return m, checkServer(m)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf(
		"> %s\n\n%s",
		m.input.View(),
		"(Ctrl+C to quit)",
	) + "\n"
}

package info

// Fetch a user's basic Charm account info

import (
	"github.com/charmbracelet/charm"
	"github.com/charmbracelet/charm/ui/common"
	"github.com/charmbracelet/tea"
	"github.com/charmbracelet/teaparty/spinner"
	te "github.com/muesli/termenv"
)

var (
	color    = te.ColorProfile().Color
	purpleBg = "#5A56E0"
	purpleFg = "#7571F9"
)

// MSG

type GotBioMsg *charm.User

// MODEL

type Model struct {
	Quit    bool // signals it's time to exit the whole application
	Err     error
	User    *charm.User
	cc      *charm.Client
	spinner spinner.Model
}

func NewModel(cc *charm.Client) Model {
	s := spinner.NewModel()
	s.Type = spinner.Dot
	s.ForegroundColor = "244"

	return Model{
		Quit:    false,
		User:    nil,
		cc:      cc,
		spinner: s,
	}
}

// UPDATE

func Update(msg tea.Msg, m Model) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			fallthrough
		case "esc":
			fallthrough
		case "ctrl+c":
			m.Quit = true
			return m, nil
		}
	case GotBioMsg:
		m.User = msg
	case tea.ErrMsg:
		// If there's an error we print the error and exit
		m.Err = msg
		m.Quit = true
		return m, nil
	case spinner.TickMsg:
		m.spinner, _ = spinner.Update(msg, m.spinner)
	}
	return m, nil
}

// VIEW

func View(m Model) string {
	if m.Err != nil {
		return "error: " + m.Err.Error()
	} else if m.User == nil {
		return spinner.View(m.spinner) + " Authenticating..."
	}
	return bioView(m.User)
}

func bioView(u *charm.User) string {
	var username string
	if u.Name != "" {
		username = u.Name
	} else {
		username = te.String("(none set)").Foreground(color("241")).String()
	}
	return common.KeyValueView(
		"Username", username,
		"Joined", u.CreatedAt.Format("02 Jan 2006"),
	)
}

// SUBSCRIPTIONS

// Tick just wraps the spinner's subscription
func Tick(model tea.Model) tea.Sub {
	m, ok := model.(Model)
	if !ok {
		// TODO: handle this error properly
		return nil
	} else if m.User != nil {
		return nil
	}

	return tea.SubMap(spinner.Sub, m.spinner)
}

// COMMANDS

func GetBio(model tea.Model) tea.Msg {
	m, ok := model.(Model)
	if !ok {
		return tea.ModelAssertionErr
	}

	user, err := m.cc.Bio()
	if err != nil {
		return tea.NewErrMsgFromErr(err)
	}

	return GotBioMsg(user)
}

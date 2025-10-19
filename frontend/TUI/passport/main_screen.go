package passport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	art "pwd-manager-tui/artistics"
	"pwd-manager-tui/auth"
	"pwd-manager-tui/util"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type mainScreenModel struct {
	Hosts        []string
	HostIndex    *int
	CurrentFocus *int
	Width        int
}

func NewMainScreenModel() *mainScreenModel {
	return &mainScreenModel{
		CurrentFocus: new(int),
		HostIndex:    new(int),
	}
}

func GenerateHosts() []string {

	request, err := http.NewRequest("GET", "https://localhost:443/pwd/getHosts", bytes.NewBuffer([]byte{}))
	if err != nil {
		return []string{}
	}

	authToken := auth.GetSavedData().AuthToken
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	response, err := auth.Client.Do(request)
	if err != nil {
		return []string{}
	}

	var buffer []byte
	if response.StatusCode == 200 {
		buffer, err = io.ReadAll(response.Body)
		if err != nil {
			return []string{}
		}
	} else {
		return []string{}
	}

	var res []string

	err = json.Unmarshal(buffer, &res)
	if err != nil {
		return []string{}
	}

	return res
}

func (model mainScreenModel) Init() tea.Cmd {
	return textinput.Blink
}

func (model mainScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+x", "q", "Q":
			return model, tea.Quit
		case "s", "S":
			err := auth.SignOut()
			if err != nil {
				return model, nil
			}
			return model, tea.Quit
		case "l", "L":
			*model.CurrentFocus = 0
		case "n", "N":
			*model.CurrentFocus = 1
		case "tab", "shift+tab":
			s := msg.String()
			if s == "tab" {
				*model.CurrentFocus++
			} else {
				*model.CurrentFocus--
			}
			if *model.CurrentFocus > 3 {
				*model.CurrentFocus = 0
			} else if *model.CurrentFocus < 0 {
				*model.CurrentFocus = 3
			}
		case "enter":
			switch *model.CurrentFocus {
			case 2:
				err := auth.SignOut()
				if err != nil {
					return model, nil
				}
				return model, tea.Quit
			case 3:
				return model, tea.Quit
			}
		case "j", "J", "down", "k", "K", "up":
			if *model.CurrentFocus == 0 {
				s := msg.String()
				switch s {
				case "j", "J", "down":
					*model.HostIndex++
				case "k", "K", "up":
					*model.HostIndex--
				}

			}
		}
	case tea.WindowSizeMsg:
		model.Width = msg.Width
	}

	return model, nil
}

func (model *mainScreenModel) ViewHosts() string {
	var builder strings.Builder

	if *model.HostIndex > len(model.Hosts)-1 {
		*model.HostIndex = 0
	} else if *model.HostIndex < 0 {
		*model.HostIndex = len(model.Hosts) - 1
	}

	for i, host := range model.Hosts {
		if *model.HostIndex == i {
			builder.WriteString(strings.Repeat("-", model.Width/2))
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("> %s", util.RightPad(host, 20)))
			builder.WriteString(art.GreenStyle.Render("\t[ View (v) ]"))
			builder.WriteString(art.YellowStyle.Render("\t[ Copy to clipboard (c) ]"))
			builder.WriteString(art.DangerStyle.Render("\t[ Remove (x) ]"))
			builder.WriteString("\n")
			builder.WriteString(strings.Repeat("-", model.Width/2))
		} else {
			builder.WriteString(fmt.Sprintf("  %s", host))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func (model mainScreenModel) View() string {
	var builder strings.Builder
	builder.WriteString(art.FocusedStyle.Render(art.LoadTitle()))
	builder.WriteString("\n")

	userInfo := auth.GetSavedData()

	userInfoString := fmt.Sprintf("%s %s [ %s ]", userInfo.Name, userInfo.Surname, userInfo.Email)

	builder.WriteString(art.FocusedStyle.Render(userInfoString))
	builder.WriteString("\n")

	builder.WriteString(art.FocusedStyle.Render(strings.Repeat("-", model.Width)))

	signOutButton := &art.BlurredSignOutButton
	quitButton := &art.BlurredQuitButton
	listButton := &art.BlurredListButton
	newPasswordButton := &art.BlurredNewPasswordButton

	pageString := ""

	switch *model.CurrentFocus {
	case 0:
		listButton = &art.FocusedListButton
		if len(model.Hosts) == 0 {
			model.Hosts = GenerateHosts()
		}
		pageString = model.ViewHosts()
	case 1:
		newPasswordButton = &art.FocusedNewPasswordButton
		pageString = model.NewPasswordView()
	case 2:
		signOutButton = &art.FocusedSignOutButton
		pageString = "Press ENTER to sign out"
	case 3:
		quitButton = &art.FocusedQuitButton
		pageString = "Press ENTER to quit passport"
	}

	fmt.Fprintf(&builder, "\n[l] %s\t\t[n] %s\t\t[s] %s\t\t[q] %s\n",
		*listButton,
		*newPasswordButton,
		*signOutButton,
		*quitButton)
	builder.WriteString(art.FocusedStyle.Render(strings.Repeat("-", model.Width)))
	builder.WriteString("\n")
	builder.WriteString(pageString)

	return builder.String()
}

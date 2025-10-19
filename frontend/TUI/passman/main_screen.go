package passman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	art "pwd-manager-tui/artistics"
	"pwd-manager-tui/auth"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type mainScreenModel struct {
	Hosts        []string
	CurrentFocus *int
	Width        int
}

func NewMainScreenModel() *mainScreenModel {
	return &mainScreenModel{
		CurrentFocus: new(int),
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
		case "tab", "shift+tab":
			s := msg.String()
			if s == "tab" {
				*model.CurrentFocus++
			} else {
				*model.CurrentFocus--
			}

			if *model.CurrentFocus > 1 {
				*model.CurrentFocus = 0
			} else if *model.CurrentFocus < 0 {
				*model.CurrentFocus = 1
			}
		case "enter":
			switch *model.CurrentFocus {
			case 0:
				err := auth.SignOut()
				if err != nil {
					return model, nil
				}
				return model, tea.Quit
			case 1:
				return model, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		model.Width = msg.Width
	}

	return model, nil
}

func (model mainScreenModel) View() string {
	var builder strings.Builder
	builder.WriteString(art.FocusedStyle.Render(art.LoadTitle()))
	builder.WriteString("\n")
	
	userInfo := auth.GetSavedData()

	userInfoString := fmt.Sprintf("%s %s [ %s ]",userInfo.Name, userInfo.Surname, userInfo.Email)

	builder.WriteString(art.FocusedStyle.Render(userInfoString))
	builder.WriteString("\n")

	builder.WriteString(art.FocusedStyle.Render(strings.Repeat("-", model.Width)))

	signOutButton := &art.BlurredSignOutButton
	quitButton := &art.BlurredQuitButton
	switch *model.CurrentFocus {
	case 0:
		signOutButton = &art.FocusedSignOutButton
	case 1:
		quitButton = &art.FocusedQuitButton
	}
	
	fmt.Fprintf(&builder, "\n\n[s] %s\t\t[q] %s\n", *signOutButton, *quitButton)
		

	return builder.String()
}

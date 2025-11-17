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

type field struct {
	Submition string
	Field     textinput.Model
}

type mainScreenModel struct {
	Hosts                   []string
	HostIndex               *int
	CurrentNewPasswordField *int
	CurrentFocus            *int
	NewPasswordInputs       []field
	Width                   int
}

func NewMainScreenModel() *mainScreenModel {
	model := &mainScreenModel{
		CurrentFocus:            new(int),
		HostIndex:               new(int),
		CurrentNewPasswordField: new(int),
		NewPasswordInputs:       make([]field, 2),
	}

	var input textinput.Model
	for i := range model.NewPasswordInputs {
		input = textinput.New()
		input.Cursor.Style = art.CursorStyle
		input.CharLimit = 100

		switch i {
		case 0:
			input.Focus()
			input.PromptStyle = art.FocusedStyle
			input.TextStyle = art.FocusedStyle
		}
		model.NewPasswordInputs[i].Field = input
	}

	return model
}


type Hosts struct {
	Hosts []string `json:"hosts"`
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

	var res Hosts

	err = json.Unmarshal(buffer, &res)
	if err != nil {
		return []string{}
	}

	return res.Hosts
}

func (model mainScreenModel) Init() tea.Cmd {
	return textinput.Blink
}

func (model mainScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+x":
			return model, tea.Quit
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
			case 1:
				if *model.CurrentNewPasswordField == 2 {
					err := model.CreateNewPassword(model.NewPasswordInputs[0].Submition, model.NewPasswordInputs[1].Submition)
					if err != nil {
						return model, tea.Quit
					}
					model.NewPasswordInputs[0].Submition = ""
					model.NewPasswordInputs[1].Submition = ""
					model.NewPasswordInputs[0].Field.SetValue("")
					model.NewPasswordInputs[1].Field.SetValue("")
					*model.CurrentFocus = 0
					return model, nil
				}

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
			s := msg.String()
			switch *model.CurrentFocus {
			case 0:
				switch s {
				case "j", "J", "down":
					*model.HostIndex++
				case "k", "K", "up":
					*model.HostIndex--
				}
			case 1:
				switch s {
				case "down":
					*model.CurrentNewPasswordField++
				case "up":
					*model.CurrentNewPasswordField--
				}

				if *model.CurrentNewPasswordField > 2 {
					*model.CurrentNewPasswordField = 0
				} else if *model.CurrentNewPasswordField < 0 {
					*model.CurrentNewPasswordField = 2
				}
				cmds := make([]tea.Cmd, len(model.NewPasswordInputs))
				for i := 0; i <= len(model.NewPasswordInputs)-1; i++ {
					if i == *model.CurrentNewPasswordField {

						cmds[i] = model.NewPasswordInputs[i].Field.Focus()
						model.NewPasswordInputs[i].Field.PromptStyle = art.FocusedStyle
						model.NewPasswordInputs[i].Field.TextStyle = art.FocusedStyle
						continue
					}

					model.NewPasswordInputs[i].Field.Blur()
					model.NewPasswordInputs[i].Field.PromptStyle = art.NoStyle
					model.NewPasswordInputs[i].Field.TextStyle = art.NoStyle
				}

			}
		}
	case tea.WindowSizeMsg:
		model.Width = msg.Width
	}

	cmd := model.updateInputs(msg)

	return model, cmd
}

func (model *mainScreenModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(model.NewPasswordInputs))
	for i := range model.NewPasswordInputs {
		model.NewPasswordInputs[i].Field, cmds[i] = model.NewPasswordInputs[i].Field.Update(msg)
		model.NewPasswordInputs[i].Submition = model.NewPasswordInputs[i].Field.Value()
	}

	return tea.Batch(cmds...)
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

	fmt.Fprintf(&builder, "\n%s\t\t%s\t\t%s\t\t%s\n",
		*listButton,
		*newPasswordButton,
		*signOutButton,
		*quitButton)
	builder.WriteString(art.FocusedStyle.Render(strings.Repeat("-", model.Width)))
	builder.WriteString("\n")
	builder.WriteString(pageString)

	return builder.String()
}

package auth

import (
	"fmt"
	art "pwd-manager-tui/artistics"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type loginModel struct {
	Inputs     []field
	FocusIndex int
	CursorMode cursor.Mode
}

type field struct {
	Submition string
	Field     textinput.Model
}

func NewLoginModel() loginModel {
	loginModel := loginModel{
		Inputs: make([]field, 2),
	}

	var textInput textinput.Model
	for i := range loginModel.Inputs {
		textInput = textinput.New()
		textInput.Cursor.Style = art.CursorStyle
		textInput.CharLimit = 100

		switch i {
		case 0:
			textInput.Focus()
			textInput.PromptStyle = art.FocusedStyle
			textInput.TextStyle = art.FocusedStyle
		case 1:
			textInput.EchoMode = textinput.EchoPassword
			textInput.EchoCharacter = '*'
		}

		loginModel.Inputs[i].Field = textInput
	}

	return loginModel
}

func (model loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (model loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+x":
			return model, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			if s == "enter" && model.FocusIndex == len(model.Inputs) {
				return model, tea.Quit
			}
			if s == "up" || s == "shift+tab" {
				model.FocusIndex--
			} else {
				model.FocusIndex++
			}

			if model.FocusIndex > len(model.Inputs) {
				model.FocusIndex = 0
			} else if model.FocusIndex < 0 {
				model.FocusIndex = len(model.Inputs)
			}

			cmds := make([]tea.Cmd, len(model.Inputs))
			for i := 0; i <= len(model.Inputs)-1; i++ {
				if i == model.FocusIndex {

					cmds[i] = model.Inputs[i].Field.Focus()
					model.Inputs[i].Field.PromptStyle = art.FocusedStyle
					model.Inputs[i].Field.TextStyle = art.FocusedStyle
					continue
				}

				model.Inputs[i].Field.Blur()
				model.Inputs[i].Field.PromptStyle = art.NoStyle
				model.Inputs[i].Field.TextStyle = art.NoStyle
			}

			return model, tea.Batch(cmds...)
		}
	}

	cmd := model.updateInputs(msg)
	return model, cmd
}

func (model *loginModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(model.Inputs))
	for i := range model.Inputs {
		model.Inputs[i].Field, cmds[i] = model.Inputs[i].Field.Update(msg)
		model.Inputs[i].Submition = model.Inputs[i].Field.Value()
	}

	return tea.Batch(cmds...)
}

func (model loginModel) View() string {
	var builder strings.Builder

	builder.WriteString(art.FocusedStyle.Render(art.LoadTitle()))

	builder.WriteString("\nEmail:\n")
	builder.WriteString(model.Inputs[0].Field.View())
	builder.WriteString("\nPassword:\n")
	builder.WriteString(model.Inputs[1].Field.View())

	button := &art.BlurredLoginButton
	if model.FocusIndex == len(model.Inputs) {
		button = &art.FocusedLoginButton
	}
	fmt.Fprintf(&builder, "\n\n%s\n\n", *button)

	return builder.String()
}

func (model *loginModel) GetValues() []string {
	result := []string{}
	result = append(result, model.Inputs[0].Submition)
	result = append(result, model.Inputs[1].Submition)

	return result
}

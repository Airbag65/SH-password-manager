package auth

import (
	"fmt"
	art "pwd-manager-tui/artistics"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)



type signUpModel struct {
	Inputs     []field
	FocusIndex int
	CursorMode cursor.Mode
}

func NewSignUpModel() signUpModel {
	signUpModel := signUpModel{
		Inputs: make([]field, 5),
	}

	
	var textInput textinput.Model
	for i := range signUpModel.Inputs {
		textInput = textinput.New()
		textInput.Cursor.Style = cursorStyle
		textInput.CharLimit = 100

		switch i {
		case 0:
			// Email
			textInput.Focus()
			textInput.PromptStyle = focusedStyle
			textInput.TextStyle = focusedStyle
		case 1:
			// Name
			// textInput.PromptStyle = blurredStyle
			textInput.TextStyle = focusedStyle
		case 2:
			// Surname
			// textInput.PromptStyle = blurredStyle
			textInput.TextStyle = focusedStyle
		case 3:
			// Password
			textInput.EchoMode = textinput.EchoPassword
			textInput.EchoCharacter = '*'
		case 4:
			// Confirm password
			textInput.EchoMode = textinput.EchoPassword
			textInput.EchoCharacter = '*'
		}
		signUpModel.Inputs[i].Field = textInput
	}
	return signUpModel
}

func (model signUpModel) Init() tea.Cmd {
	return textinput.Blink
}


func (model signUpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
					model.Inputs[i].Field.PromptStyle = focusedStyle
					model.Inputs[i].Field.TextStyle = focusedStyle
					continue
				}

				model.Inputs[i].Field.Blur()
				model.Inputs[i].Field.PromptStyle = noStyle
				model.Inputs[i].Field.TextStyle = noStyle
			}

			return model, tea.Batch(cmds...)
		}
	}

	cmd := model.updateInputs(msg)
	return model, cmd
}

func (model *signUpModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(model.Inputs))
	for i := range model.Inputs {
		model.Inputs[i].Field, cmds[i] = model.Inputs[i].Field.Update(msg)
		model.Inputs[i].Submition = model.Inputs[i].Field.Value()
	}

	return tea.Batch(cmds...)
}

func (model signUpModel) View() string {
	var builder strings.Builder

	builder.WriteString(focusedStyle.Render(art.LoadTitle()))

	builder.WriteString("\nEmail:\n")
	builder.WriteString(model.Inputs[0].Field.View())
	builder.WriteString("\nName:\n")
	builder.WriteString(model.Inputs[1].Field.View())
	builder.WriteString("\nSurname:\n")
	builder.WriteString(model.Inputs[2].Field.View())
	builder.WriteString("\nPassword:\n")
	builder.WriteString(model.Inputs[3].Field.View())
	builder.WriteString("\nConfirm password:\n")
	builder.WriteString(model.Inputs[4].Field.View())

	button := &blurredSignUpButton
	if model.FocusIndex == len(model.Inputs) {
		button = &focusedSignUpButton
	}
	fmt.Fprintf(&builder, "\n\n%s\n\n", *button)
	
	return builder.String()
}

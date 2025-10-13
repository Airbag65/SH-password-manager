package auth

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#05a317"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()

	focusedButton = focusedStyle.Render("[ Login ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Login"))
)

const titleString = `
 /$$$$$$$                              /$$      /$$                    
| $$__  $$                            | $$$    /$$$                    
| $$  \ $$ /$$$$$$   /$$$$$$$ /$$$$$$$| $$$$  /$$$$  /$$$$$$  /$$$$$$$ 
| $$$$$$$/|____  $$ /$$_____//$$_____/| $$ $$/$$ $$ |____  $$| $$__  $$
| $$____/  /$$$$$$$|  $$$$$$|  $$$$$$ | $$  $$$| $$  /$$$$$$$| $$  \ $$
| $$      /$$__  $$ \____  $$\____  $$| $$\  $ | $$ /$$__  $$| $$  | $$
| $$     |  $$$$$$$ /$$$$$$$//$$$$$$$/| $$ \/  | $$|  $$$$$$$| $$  | $$
|__/      \_______/|_______/|_______/ |__/     |__/ \_______/|__/  |__/
`

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
		textInput.Cursor.Style = cursorStyle
		textInput.CharLimit = 100

		switch i {
		case 0:
			// textInput.Placeholder = "Email"
			textInput.Focus()
			textInput.PromptStyle = focusedStyle
			textInput.TextStyle = focusedStyle
		case 1:
			// textInput.Placeholder = "Password"
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

	builder.WriteString(focusedStyle.Render(titleString))

	builder.WriteString("\nEmail:\n")
	builder.WriteString(model.Inputs[0].Field.View())
	builder.WriteString("\nPassword:\n")
	builder.WriteString(model.Inputs[1].Field.View())

	button := &blurredButton
	if model.FocusIndex == len(model.Inputs) {
		button = &focusedButton
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

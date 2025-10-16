package auth

import (
	"fmt"
	art "pwd-manager-tui/artistics"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type startScreenModel struct {
	FocusIndex *int
}

func NewStartScreenModel(startValue *int) startScreenModel {
	return startScreenModel{
		FocusIndex: startValue,
	}
}

func (model startScreenModel) Init() tea.Cmd {
	return textinput.Blink
}

func (model startScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+x":
			return model, tea.Quit
		case "tab", "shift+tab", "right", "left", "enter", "1", "2":
			s := msg.String()
			if s == "enter" {
				return model, tea.Quit
			}
			if s == "tab" || s == "right" {
				*model.FocusIndex++
			}

			if s == "shift+tab" || s == "left" {
				*model.FocusIndex--
			}

			if *model.FocusIndex > 1 {
				*model.FocusIndex = 0
			} else if *model.FocusIndex < 0 {
				*model.FocusIndex = 1
			}
			return model, nil
		}
	}
	return model, nil
}

func (model *startScreenModel) GetValue() int {
	return *model.FocusIndex
}

func (model startScreenModel) View() string {
	var builder strings.Builder

	builder.WriteString(focusedStyle.Render(art.LoadTitle()))

	builder.WriteString("\n\n")
	if *model.FocusIndex == 0 {
		fmt.Fprintf(&builder, "\t%s\t", focusedLoginButton)
		fmt.Fprintf(&builder, "\t%s\t", blurredSignUpButton)
	} else {
		fmt.Fprintf(&builder, "\t%s\t", blurredLoginButton)
		fmt.Fprintf(&builder, "\t%s\t", focusedSignUpButton)
	}

	return builder.String()
}

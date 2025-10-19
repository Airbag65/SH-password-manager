package artistics

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#21aaff"))
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	DangerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	GreenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#007800"))
	YellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cad902"))

	CursorStyle = FocusedStyle
	NoStyle     = lipgloss.NewStyle()

	FocusedLoginButton = FocusedStyle.Render("[ Login ]")
	BlurredLoginButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Login"))

	FocusedSignUpButton = FocusedStyle.Render("[ Sign up ]")
	BlurredSignUpButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Sign up"))

	FocusedSignOutButton = FocusedStyle.Render("[ Sign out ]")
	BlurredSignOutButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Sign out"))

	FocusedQuitButton = FocusedStyle.Render("[ Quit ]")
	BlurredQuitButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Quit"))

	FocusedListButton = FocusedStyle.Render("[ List ]")
	BlurredListButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("List"))

	FocusedNewPasswordButton = FocusedStyle.Render("[ New password ]")
	BlurredNewPasswordButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("New password"))
)

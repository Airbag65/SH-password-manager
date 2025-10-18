package passman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pwd-manager-tui/auth"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type mainScreenModel struct {
	Hosts        []string
	CurrentFocud int
}

func NewMainScreenModel() *mainScreenModel {
	return &mainScreenModel{}
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



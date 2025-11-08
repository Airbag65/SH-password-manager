package passport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pwd-manager-tui/artistics"
	"pwd-manager-tui/auth"
	"pwd-manager-tui/enc"
	"slices"
	"strings"
)

type createPasswordRequest struct {
	HostName string `json:"host_name"`
	Password string `json:"password"`
}

func (model *mainScreenModel) CreateNewPassword(hostName, password string) error {
	if slices.Contains(model.Hosts, hostName) {
		return fmt.Errorf("password for host: '%s' already exists", hostName)
	}

	rsaPemString, err := enc.PEMFileToString()
	if err != nil {
		return err
	}

	publicKey := enc.PemStringToPublicKey(rsaPemString)

	encPassword, err := enc.Encrypt(password, publicKey)
	if err != nil {
		return err
	}

	newPasswordReq := createPasswordRequest{
		HostName: hostName,
		Password: string(encPassword),
	}

	reqBody, err := json.Marshal(newPasswordReq)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", "https://localhost:443/pwd/new", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	authToken := auth.GetSavedData().AuthToken

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	
	response, err := auth.Client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("Failed to create new password, statusCode was: %d\n", response.StatusCode)
	}

	return nil
}

func (model *mainScreenModel) NewPasswordView() string {
	var builder strings.Builder
	

	builder.WriteString("Add new password\n")
	builder.WriteString("\nHost:\n")
	builder.WriteString(model.NewPasswordInputs[0].Field.View())
	builder.WriteString("\nPassword:\n")
	builder.WriteString(model.NewPasswordInputs[1].Field.View())
	button := &artistics.BlurredNewPasswordButton
	if *model.CurrentNewPasswordField == 2 {
		button = &artistics.FocusedNewPasswordButton
	}
	fmt.Fprintf(&builder, "\n\n%s\n\n", *button)

	return builder.String()
}

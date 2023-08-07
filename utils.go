package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func getOrgID(orgName string) (string, error) {
	result, err := cliConnection.CliCommandWithoutTerminalOutput("curl", "/v3/organizations?names="+orgName)
	if err != nil {
		return "", err
	}
	var organizations struct {
		Resources []Organization `json:"resources"`
	}
	err = json.Unmarshal([]byte(strings.Join(result, "\n")), &organizations)
	if err != nil {
		return "", err
	}
	if len(organizations.Resources) == 0 {
		return "", fmt.Errorf("org %s not found", orgName)
	}
	return organizations.Resources[0].GUID, nil
}

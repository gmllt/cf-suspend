package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/orange-cloudfoundry/cf-security-entitlement/plugin/messages"
)

type SuspendedState struct {
	Suspended bool `json:"suspended"`
}

type StateOptions struct {
	Org string `positional-arg-name:"ORG" required:"true" description:"Organization to resume"`
}

type StateCommand struct {
	API           string        `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	ResumeOptions ResumeOptions `required:"2" positional-args:"true"`
}

var stateCommand StateCommand

func (c *StateCommand) Execute(_ []string) error {
	username, err := cliConnection.Username()
	if err != nil {
		return err
	}

	orgShow := c.ResumeOptions.Org
	if orgShow != "" {
		orgShow = fmt.Sprint(messages.C.Cyan(orgShow))
	}
	_, _ = messages.Printf("Checking suspended state of org %s as %s\n", orgShow, messages.C.Cyan(username))

	orgID, err := getOrgID(c.ResumeOptions.Org)
	if err != nil {
		return err
	}

	state, err := isOrgSuspended(orgID)
	if err != nil {
		return err
	}
	strState := messages.C.Red("suspended")
	if !state {
		strState = messages.C.Green("not suspended")
	}
	_, _ = messages.Printf("Organization %s is %s\n", messages.C.Cyan(c.ResumeOptions.Org), strState)
	return nil
}

func isOrgSuspended(orgGUID string) (bool, error) {
	result, err := cliConnection.CliCommandWithoutTerminalOutput("curl", "-X", "GET", "/v3/organizations/"+orgGUID)
	if err != nil {
		return false, err
	}
	var suspendedState SuspendedState
	err = json.Unmarshal([]byte(strings.Join(result, "")), &suspendedState)
	if err != nil {
		return false, err
	}

	return suspendedState.Suspended, nil
}

func init() {
	desc := `Get suspended state of an organization.`
	_, err := parser.AddCommand(
		"is-org-suspended",
		desc,
		desc,
		&stateCommand)
	if err != nil {
		panic(err)
	}
}

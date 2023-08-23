package main

import (
	"fmt"

	"github.com/orange-cloudfoundry/cf-security-entitlement/plugin/messages"
)

type ResumeOptions struct {
	Org string `positional-arg-name:"ORG" required:"true" description:"Organization to resume"`
}

type ResumeCommand struct {
	API           string        `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	ResumeOptions ResumeOptions `required:"2" positional-args:"true"`
}

var resumeCommand ResumeCommand

func (c *ResumeCommand) Execute(_ []string) error {
	username, err := cliConnection.Username()
	if err != nil {
		return err
	}

	orgShow := c.ResumeOptions.Org
	if orgShow != "" {
		orgShow = fmt.Sprint(messages.C.Cyan(orgShow))
	}
	_, _ = messages.Printf("Resuming org %s as %s\n", orgShow, messages.C.Cyan(username))

	orgID, err := getOrgID(c.ResumeOptions.Org)
	if err != nil {
		return err
	}

	err = resumeOrg(orgID)
	if err != nil {
		return err
	}
	_, _ = messages.Println(messages.C.Green("OK\n"))
	return nil
}

func resumeOrg(orgGUID string) error {
	_, err := cliConnection.CliCommandWithoutTerminalOutput("curl", "-X", "PATCH", "/v3/organizations/"+orgGUID, "-d", `{"suspended":false}`)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	desc := `Resume a organization.`
	_, err := parser.AddCommand(
		"resume-org",
		desc,
		desc,
		&resumeCommand)
	if err != nil {
		panic(err)
	}
}

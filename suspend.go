package main

import (
	"fmt"

	"github.com/orange-cloudfoundry/cf-security-entitlement/plugin/messages"
)

type SuspendOptions struct {
	Org string `positional-arg-name:"ORG" required:"true" description:"Organization to suspend"`
}

type SuspendCommand struct {
	API            string         `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
	SuspendOptions SuspendOptions `required:"2" positional-args:"true"`
}

func (c *SuspendCommand) Execute(_ []string) error {
	username, err := cliConnection.Username()
	if err != nil {
		return err
	}

	orgShow := c.SuspendOptions.Org
	if orgShow != "" {
		orgShow = "/" + fmt.Sprint(messages.C.Cyan(orgShow))
	}
	_, _ = messages.Printf("Suspending org %s as %s\n", orgShow, messages.C.Cyan(username))

	orgID, err := getOrgID(c.SuspendOptions.Org)
	if err != nil {
		return err
	}

	err = suspendOrg(orgID)
	if err != nil {
		return err
	}
	_, _ = messages.Println(messages.C.Green("OK\n"))
	return nil
}

func suspendOrg(orgGUID string) error {
	_, err := cliConnection.CliCommandWithoutTerminalOutput("curl", "-X", "PATCH", "/v3/organizations/"+orgGUID, "-d", `{"suspend":true}`)
	if err != nil {
		return err
	}

	return nil
}

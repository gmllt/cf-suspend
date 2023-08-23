package main

import (
	"encoding/json"
	"github.com/orange-cloudfoundry/cf-security-entitlement/plugin/messages"
	"strconv"
	"strings"
)

type SuspendedListResult struct {
	Resources []SuspendedOrg `json:"resources"`
}

type SuspendedOrg struct {
	SuspendedState
	OrgName string `json:"name"`
}

type ListOptions struct {
}

type ListCommand struct {
	API string `short:"a" long:"api" description:"API endpoint (e.g. https://api.example.com)"`
}

var listCommand ListCommand

func (c *ListCommand) Execute(_ []string) error {
	username, err := cliConnection.Username()
	if err != nil {
		return err
	}

	_, _ = messages.Printf("Getting suspended organizations as %s\n", messages.C.Cyan(username))

	orgs, err := listSuspendedOrgs()
	if err != nil {
		return err
	}

	if len(orgs) == 0 {
		_, _ = messages.Println("No suspended organizations")
	} else {
		_, _ = messages.Printf("%s suspended organizations:\n", messages.C.Cyan(strconv.Itoa(len(orgs))))
		for _, org := range orgs {
			_, _ = messages.Printf(" - %s\n", messages.C.Red(org))
		}
	}

	return nil
}

func listSuspendedOrgs() ([]string, error) {
	result, err := cliConnection.CliCommandWithoutTerminalOutput("curl", "-X", "GET", "/v3/organizations?per_page=5000")
	if err != nil {
		return nil, err
	}
	var suspendedList SuspendedListResult
	err = json.Unmarshal([]byte(strings.Join(result, "")), &suspendedList)
	if err != nil {
		return nil, err
	}
	var suspendedOrgs []string
	for _, org := range suspendedList.Resources {
		if org.Suspended {
			suspendedOrgs = append(suspendedOrgs, org.OrgName)
		}
	}
	return suspendedOrgs, nil
}

func init() {
	desc := `List suspended organizations.`
	_, err := parser.AddCommand(
		"list-suspended-orgs",
		desc,
		desc,
		&listCommand)
	if err != nil {
		panic(err)
	}
}

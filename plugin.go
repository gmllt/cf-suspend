package main

import (
	"errors"
	"fmt"
	"io"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/jessevdk/go-flags"
	"github.com/orange-cloudfoundry/cf-security-entitlement/plugin/messages"
)

type Options struct {
}

var pluginVersion = "0.1.0"
var options Options
var parser = flags.NewParser(&options, flags.HelpFlag|flags.PassDoubleDash|flags.IgnoreUnknown)
var cliConnection plugin.CliConnection

func Parse(args []string) error {
	_, err := parser.ParseArgs(args)
	if err != nil {
		var errFlag *flags.Error
		if errors.As(err, &errFlag) && errFlag.Type == flags.ErrCommandRequired {
			return nil
		}
		if errors.As(err, &errFlag) && errFlag.Type == flags.ErrHelp {
			messages.Errorf("Error parsing arguments: %s", err)
			return nil
		}
		return err
	}

	return nil
}

type SuspendPlugin struct {
	Connection plugin.CliConnection
	Out        io.Writer
}

func (c *SuspendPlugin) GetMetadata() plugin.PluginMetadata {
	var major, minor, build int
	fmt.Sscanf(pluginVersion, "%d.%d.%d", &major, &minor, &build)
	return plugin.PluginMetadata{
		Name: "SuspendPlugin",
		Version: plugin.VersionType{
			Major: major,
			Minor: minor,
			Build: build,
		},
		Commands: []plugin.Command{
			{
				Name:     "suspend-org",
				HelpText: "Suspend an organization.",
				UsageDetails: plugin.Usage{
					Usage: "cf suspend-org ORG_NAME",
				},
			},
			{
				Name:     "resume-org",
				HelpText: "Resume an organization.",
				UsageDetails: plugin.Usage{
					Usage: "cf resume-org ORG_NAME",
				},
			},
			{
				Name:     "is-org-suspended",
				HelpText: "Check if an organization is suspended.",
				UsageDetails: plugin.Usage{
					Usage: "cf is-org-suspended ORG_NAME",
				},
			},
		},
	}
}

func (c *SuspendPlugin) Run(cc plugin.CliConnection, args []string) {
	cliConnection = cc

	action := args[0]
	if action == "CLI-MESSAGE-UNINSTALL" {
		return
	}

	err := Parse(args)
	if err != nil {
		messages.Fatal(err.Error())
	}
}

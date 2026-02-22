package cli

import (
	"github.com/cgons/hdinfo/internal/cli/commands"
	urfcli "github.com/urfave/cli/v3"
)

func Register() *urfcli.Command {
	return &urfcli.Command{
		Name:  "hdinfo",
		Usage: "Displays hard drive details: size, model, mounts etc...",
		UsageText: `Note: hdinfo depends on tools like lsblk, hdparm and smartmontools.
		  (Please ensure required tools are installed system-wide)
		  (see https://github.com/cgons/hdinfo#deps for more details)

Use the commands below to get started.`,
		Commands: []*urfcli.Command{
			commands.DisksCommand(),
			commands.MountsCommand(),
		},
	}
}

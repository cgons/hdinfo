package cli

import (
	"context"

	"github.com/cgons/hdinfo/internal/cli/commands"
	"github.com/fatih/color"
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
		Flags: []urfcli.Flag{
			&urfcli.BoolFlag{
				Name:  "no-color",
				Usage: "Disable colorized output",
			},
		},
		Before: func(ctx context.Context, cmd *urfcli.Command) (context.Context, error) {
			if cmd.Bool("no-color") {
				color.NoColor = true
			}
			return ctx, nil
		},
		Commands: []*urfcli.Command{
			commands.DisksCommand(),
			commands.MountsCommand(),
		},
	}
}

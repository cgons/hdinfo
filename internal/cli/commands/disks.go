package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cgons/hdinfo/internal/lib/utils"
	"github.com/cgons/hdinfo/internal/services/diskservice"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v3"
	urfcli "github.com/urfave/cli/v3"
)

func DisksCommand() *urfcli.Command {
	return &urfcli.Command{
		Name:  "disks",
		Usage: "List all system disks and associated details",
		Before: func(ctx context.Context, cmd *urfcli.Command) (context.Context, error) {
			utils.RequireRoot()
			return ctx, nil
		},
		Action: disksCommandAction,
		Flags: []urfcli.Flag{
			&urfcli.BoolFlag{
				Name:    "silent",
				Aliases: []string{"s"},
				Usage:   "Suppress informational notes in output",
			},
			&urfcli.BoolFlag{
				Name:  "smart-data",
				Usage: "Include SMART data (temperature, power-on hours, power cycle count)",
			},
			&urfcli.BoolFlag{
				Name:  "force",
				Usage: "Wake sleeping disks to fetch SMART data (must be used with --smart-data)",
			},
		},
	}
}

func disksCommandAction(ctx context.Context, cmd *urfcli.Command) error {
	// Validate flag usage.
	// 1. --force can only be set with --smart-data
	if cmd.Bool("force") && !cmd.Bool("smart-data") {
		return cli.Exit("Error: --force cannot be used without --smart-data\n", 1)
	}

	showSmart := cmd.Bool("smart-data")
	diskDetails, err := diskservice.GetDiskDetails(diskservice.GetDiskDetailsOptions{
		GetSmartDetails:       showSmart,
		ForcePullSmartDetails: cmd.Bool("force"),
	})
	if err != nil {
		println("Error: Unable to get disk details")
		return nil
	}

	headerFormatter := color.New(color.Bold, color.Underline).SprintfFunc()
	textGreen := color.New(color.FgGreen).SprintFunc()
	textYellow := color.New(color.FgYellow).SprintFunc()

	headers := []any{
		headerFormatter("%s", "Name"),
		headerFormatter("%s", "Model"),
		headerFormatter("%s", "Capacity"),
		headerFormatter("%s", "IsSSD"),
		headerFormatter("%s", "Interface"),
		headerFormatter("%s", "State"),
	}
	if showSmart {
		headers = append(headers,
			headerFormatter("%s", "CurrentTemp"),
			headerFormatter("%s", "PowerOnHours"),
			headerFormatter("%s", "PowerCyleCount"),
		)
	}

	tbl := table.New(headers...).WithWidthFunc(func(s string) int {
		return visibleRuneCount(s)
	})

	for _, disk := range *diskDetails {
		if disk.Type != "disk" {
			continue
		}

		row := []any{
			disk.Name,
			textGreen(disk.Model),
			disk.Capacity,
			disk.IsSSD,
			disk.Interface,
			disk.State,
		}

		if showSmart {
			var tempVal string
			tempVal = strconv.Itoa(disk.SmartDetails.TempCurrent)
			if disk.SmartDetails.TempMin > 0 || disk.SmartDetails.TempMax > 0 {
				tempVal = fmt.Sprintf("%d (Min/Max %d/%d)",
					disk.SmartDetails.TempCurrent,
					disk.SmartDetails.TempMin,
					disk.SmartDetails.TempMax,
				)
			}
			row = append(row, tempVal, disk.SmartDetails.PowerOnHours, disk.SmartDetails.PowerCycleCount)
		}

		tbl.AddRow(row...)
	}

	if !cmd.Bool("silent") {
		println()
		println("hdinfo -", "github.com/cgons/hdinfo")
		println("-----------------------------------------")
		println(" - " + textYellow("--smart-data") + "  | view SMART data (temp, power-on-hours, etc...)")
		println(" - " + textYellow("--force") + "       | wake sleep/standby disks and pull SMART data")
		println(" - " + textYellow("-s") + " / " + textYellow("--silent") + " | silence these hints")
		println()
	}

	tbl.Print()
	println()

	return nil
}

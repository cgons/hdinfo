package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cgons/hdinfo/internal/lib/utils"
	"github.com/cgons/hdinfo/internal/lib/services/diskservice"
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
				Usage:   "Suppress header and informational notes in output",
			},
			&urfcli.BoolFlag{
				Name:  "smart-data",
				Usage: "Include SMART data (temperature, power-on hours, power cycle count)",
			},
			&urfcli.BoolFlag{
				Name:  "force",
				Usage: "Wake sleeping disks to fetch SMART data (must be used with --smart-data)",
			},
			&urfcli.BoolFlag{
				Name:  "no-stats",
				Usage: "Hide disk statistics totals from output",
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

	headerColor := color.New(color.Bold, color.Underline)
	headerColor.EnableColor()
	headerFormatter := headerColor.SprintfFunc()
	textBlue := color.New(color.FgBlue).SprintFunc()
	textGreen := color.New(color.FgGreen).SprintFunc()
	textYellow := color.New(color.FgYellow).SprintFunc()
	textWhite := color.New(color.FgWhite).SprintFunc()
	textFaint := color.New(color.Faint).SprintFunc()

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

	statTotals := diskservice.StatTotals{}

	for _, disk := range *diskDetails {
		if disk.Type != "disk" {
			continue
		}

		// Populate the row with data + coloured output as necessary
		row := []any{
			disk.Name,
			textBlue(disk.Model),
			disk.Capacity,
			disk.IsSSD,
			disk.Interface,
			colorizedState(disk.State, textGreen, textYellow, textWhite),
		}

		// Add SMART attribut columns if requested (via --smart-data flag)
		if showSmart {
			row = append(row, smartDataColumns(disk.SmartDetails)...)
		}

		tbl.AddRow(row...)
		diskservice.UpdateDiskStatTotals(&statTotals, disk)
	}

	// Print header section
	if !cmd.Bool("silent") {
		println()
		println(textFaint("hdinfo - ", utils.GetVersion(), " | github.com/cgons/hdinfo"))
		println(textFaint("-----------------------------------------"))
		println(textFaint("--help - to see all display options"))
		println()
	}

	// Print table
	tbl.Print()
	println()

	// Print footer section
	if !cmd.Bool("no-stats") {
		fmt.Printf(
			"Totals: %d Disks | %s | %s\n",
			statTotals.Total,
			textGreen(fmt.Sprintf("%d Active", statTotals.Active)),
			textYellow(fmt.Sprintf("%d Standby", statTotals.Standby)),
		)
		println()
	}

	return nil
}

// Helper Functions
// ----------------

func smartDataColumns(smart diskservice.SmartDetails) []any {
	tempVal := strconv.Itoa(smart.TempCurrent)
	if smart.TempMin > 0 || smart.TempMax > 0 {
		tempVal = fmt.Sprintf("%d (Min/Max %d/%d)", smart.TempCurrent, smart.TempMin, smart.TempMax)
	}
	return []any{tempVal, smart.PowerOnHours, smart.PowerCycleCount}
}

func colorizedState(state string, green, yellow, white func(a ...interface{}) string) string {
	switch state {
	case "active/idle":
		return green(state)
	case "standby", "sleeping":
		return yellow(state)
	case "unknown":
		return white(state)
	default:
		return state
	}
}

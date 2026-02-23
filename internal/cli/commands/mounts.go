package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cgons/hdinfo/internal/lib/utils"
	"github.com/cgons/hdinfo/internal/services/diskservice"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	urfcli "github.com/urfave/cli/v3"
)

func MountsCommand() *urfcli.Command {
	return &urfcli.Command{
		Name:  "mounts",
		Usage: "List disk partitions and mount points",
		Before: func(ctx context.Context, cmd *urfcli.Command) (context.Context, error) {
			utils.RequireRoot()
			return ctx, nil
		},
		Action: mountsCommandAction,
	}
}

func mountsCommandAction(ctx context.Context, cmd *urfcli.Command) error {
	diskDetails, err := diskservice.GetDiskDetails(diskservice.GetDiskDetailsOptions{})
	if err != nil {
		println("Error: Unable to get disk details")
		return nil
	}

	headerColor := color.New(color.Bold, color.Underline)
	headerColor.EnableColor() // force color/styling for header row
	headerFormatter := headerColor.SprintfFunc()
	modelGreen := color.New(color.FgGreen).SprintFunc()

	tbl := table.New(
		headerFormatter("%s", "Name"),
		headerFormatter("%s", "Model"),
		headerFormatter("%s", "MountPoint"),
		headerFormatter("%s", "Capacity"),
		headerFormatter("%s", "UsedSpace"),
		headerFormatter("%s", "FreeSpace"),
		headerFormatter("%s", "Used/Free"),
	).WithWidthFunc(func(s string) int {
		return visibleRuneCount(s)
	})

	parentModels := map[string]string{}
	for _, disk := range *diskDetails {
		if disk.Type == "disk" {
			parentModels[disk.Name] = disk.Model
		}
	}

	for _, disk := range *diskDetails {
		if disk.Type != "part" {
			continue
		}
		model := disk.Model
		if parentModel, ok := parentModels[disk.ParentName]; ok {
			model = parentModel
		}
		tbl.AddRow(
			disk.Name,
			modelGreen(model),
			disk.MountPoint,
			disk.Capacity,
			disk.UsedSpace,
			disk.FreeSpace,
			usedFreePercentage(disk.FreePercentage),
		)
	}
	tbl.Print()
	println()

	return nil
}

// usedFreePercentage formats a used% string (e.g. "5%") as "5% / 95%".
func usedFreePercentage(usedPct string) string {
	trimmed := strings.TrimSuffix(usedPct, "%")
	used, err := strconv.Atoi(trimmed)
	if err != nil || usedPct == "" {
		return usedPct
	}
	return fmt.Sprintf("%d%% / %d%%", used, 100-used)
}

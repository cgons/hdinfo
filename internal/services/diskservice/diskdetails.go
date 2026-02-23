package diskservice

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/cgons/hdinfo/internal/lib"
	"github.com/cgons/hdinfo/internal/lib/utils"
)

type DiskDetails struct {
	Name           string // the linx name (eg. sda, sdb)
	ParentName     string // parent block device name (eg. sda for sda1)
	Model          string
	Type           string // is this a "disk" or "part" (partition)?
	Capacity       string // capacity string as reported by lsblk (eg. 200G, 2T)
	IsSSD          bool
	Interface      string // connecting interface (eg. usb, sata, nvme)
	State          string // is the drive active or in standby?
	MountPoint     string
	UsedSpace      string
	FreeSpace      string
	FreePercentage string // amount of free space as a percentage
	SmartDetails   SmartDetails
}

var lsblkPairPattern = regexp.MustCompile(`([A-Z%]+)="([^"]*)"`)
var hdparmDriveStatePattern = regexp.MustCompile(`(?m)drive state is:\s*(.+)\s*$`)

type GetDiskDetailsOptions struct {
	GetSmartDetails       bool
	ForcePullSmartDetails bool
}

// GetDiskDetails returns details of disks and mounts.
func GetDiskDetails(opts GetDiskDetailsOptions) (*[]DiskDetails, error) {
	diskDetails, err := processRawDiskDetails(opts)
	if err != nil {
		return &[]DiskDetails{}, err
	}

	return diskDetails, nil
}

type StatTotals struct {
	Total   int
	Active  int
	Standby int
	Unknown int
}

// UpdateDiskStatTotals updates the given StatTotals based on a single DiskDetails entry.
// It should be called from within a loop iterating over disk details.
func UpdateDiskStatTotals(totals *StatTotals, d DiskDetails) {
	if d.Type != "disk" {
		return
	}
	totals.Total++
	switch d.State {
	case "active/idle":
		totals.Active++
	case "standby", "sleeping":
		totals.Standby++
	case "unknown":
		totals.Unknown++
	}
}

// Utility Functions
// ----------------------------------------------
func getRawDiskDetails() *bytes.Buffer {
	return utils.GetCommandOutput("lsblk", "-nP", "--sort", "NAME", "-o", "NAME,PKNAME,TYPE,MODEL,SIZE,ROTA,TRAN,MOUNTPOINT,FSUSED,FSAVAIL,FSUSE%")
}

func processRawDiskDetails(opts GetDiskDetailsOptions) (*[]DiskDetails, error) {
	rawDetails := getRawDiskDetails()

	var diskDetails []DiskDetails
	var rawMappings []map[string]string
	scanner := bufio.NewScanner(rawDetails)
	for scanner.Scan() {
		// Example raw output (single line):
		// NAME="sdc" PKNAME="" TYPE="disk" MODEL="ST4000DM004-2CV104" SIZE="3.6T"
		// ROTA="1" TRAN="usb" MOUNTPOINT="/boot/efi" FSUSED="19G" FSAVAIL="185.4G" FSUSE%="9%"
		// --
		// Match will resul in:
		// 10 => [
		//    0 => "FSUSE%="""
		//    1 => "FSUSE%" <-- KEY
		//    2 => "" <-- VALUE
		// ]
		matches := lsblkPairPattern.FindAllStringSubmatch(scanner.Text(), -1)

		// Extract K,V
		mappedLine := make(map[string]string)
		for _, match := range matches {
			mappedLine[match[1]] = match[2]
		}
		rawMappings = append(rawMappings, mappedLine)
	}

	// Loop over mappings and assign
	for _, mapping := range rawMappings {
		if utils.GetOrDefault(mapping, "NAME", "") != "" {
			diskType := utils.GetOrDefault(mapping, "TYPE", "")

			disk := DiskDetails{
				Name:           utils.GetOrDefault(mapping, "NAME", ""),
				ParentName:     utils.GetOrDefault(mapping, "PKNAME", ""),
				Type:           diskType,
				Model:          utils.GetOrDefault(mapping, "MODEL", ""),
				Capacity:       utils.GetOrDefault(mapping, "SIZE", ""),
				IsSSD:          utils.GetOrDefault(mapping, "ROTA", "") == "0",
				Interface:      utils.GetOrDefault(mapping, "TRAN", ""),
				State:          "",
				MountPoint:     utils.GetOrDefault(mapping, "MOUNTPOINT", ""),
				UsedSpace:      utils.GetOrDefault(mapping, "FSUSED", ""),
				FreeSpace:      utils.GetOrDefault(mapping, "FSAVAIL", ""),
				FreePercentage: utils.GetOrDefault(mapping, "FSUSE%", ""),
			}

			if diskType == "disk" {
				// Get disk state status (active or idle)
				disk.State = getDiskState(mapping["NAME"])

				// Get disk smart details (temp, powerontime, etc...)
				// Only probe if the disk is active/idle to avoid waking standby disks,
				// unless --force is set.
				if opts.GetSmartDetails && (disk.State == "active/idle" || opts.ForcePullSmartDetails) {
					disk.SmartDetails = GetSmartDetails(disk.Name)
				}
			}
			diskDetails = append(diskDetails, disk)
		}
	}

	if err := scanner.Err(); err != nil {
		return &[]DiskDetails{}, fmt.Errorf("failed to parse disk details: %w", err)
	}

	return &diskDetails, nil
}

func getDiskState(device string) string {
	// device - represents the disk "name" (eg. sda)
	logger := lib.GetLogger()
	device = strings.TrimSpace(device)
	if device == "" {
		return ""
	}

	if !strings.HasPrefix(device, "/dev/") {
		device = "/dev/" + device
	}

	out, err := exec.Command("hdparm", "-C", device).CombinedOutput()
	if err != nil {
		logger.Warnln("Unable to run 'hdparm' ", err)
		return ""
	}

	match := hdparmDriveStatePattern.FindStringSubmatch(string(out))
	if len(match) < 2 {
		return ""
	}

	return strings.TrimSpace(match[1])
}

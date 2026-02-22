package diskservice

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"github.com/cgons/hdinfo/internal/lib/utils"
	"github.com/tidwall/gjson"
)

type SmartDetails struct {
	DiskName        string
	PowerCycleCount int `json:"power_cycle_count"`
	PowerOnHours    int `json:"power_on_time.hours"`
	TempCurrent     int `json:"temperature.current"`
	TempMin         int
	TempMax         int
}

var tempMinMaxPattern = regexp.MustCompile(`Min/Max\s+(\d+)/(\d+)`)

func GetSmartDetails(diskName string) SmartDetails {
	rawDetails := getRawSmartAttributes(diskName)
	raw := rawDetails.String()

	tempMin, tempMax := parseTempMinMax(raw)

	return SmartDetails{
		DiskName:        diskName,
		PowerCycleCount: int(gjson.Get(raw, "power_cycle_count").Int()),
		PowerOnHours:    int(gjson.Get(raw, "power_on_time.hours").Int()),
		TempCurrent:     int(gjson.Get(raw, "temperature.current").Int()),
		TempMin:         tempMin,
		TempMax:         tempMax,
	}
}

func getRawSmartAttributes(diskName string) *bytes.Buffer {
	// diskName --> sda, sdb, etc...
	diskPath := fmt.Sprintf("/dev/%s", diskName)
	return utils.GetCommandOutput("smartctl", "-A", "--json", diskPath)
}

// parseTempMinMax extracts Min/Max values from the Temperature_Celsius SMART attribute (id 194)
// Example raw string: 33 (Min/Max 15/47)
func parseTempMinMax(raw string) (int, int) {
	tempRaw := gjson.Get(raw, `ata_smart_attributes.table.#(id==194).raw.string`).String()
	match := tempMinMaxPattern.FindStringSubmatch(tempRaw)
	if len(match) < 3 {
		return 0, 0
	}
	min, _ := strconv.Atoi(match[1])
	max, _ := strconv.Atoi(match[2])
	return min, max
}

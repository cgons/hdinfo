<div align="center">

# hdinfo

A command-line utility for Linux that displays hard drive details including disk information, capacity, model, mount points, and SMART data (temperature, power-on hours, power cycle count).
</div>

## Features

- List all system disks with model, capacity, interface type, and state
- View partition mount points with used/free space
- Retrieve SMART data: temperature, power-on hours, power cycle count
- Respects disk sleep states (won't wake sleeping disks unless forced)

## <a name="deps"></a>Dependencies

`hdinfo` relies on the following system utilities that must be installed:

| Tool | Package | Purpose |
|------|---------|---------|
| `lsblk` | `util-linux` | List block devices |
| `hdparm` | `hdparm` | Query disk state (active/standby) |
| `smartctl` | `smartmontools` | Read SMART attributes |

```bash
# Debian/Ubuntu
sudo apt install util-linux hdparm smartmontools

# Arch Linux
sudo pacman -S util-linux hdparm smartmontools

# Fedora/RHEL
sudo dnf install util-linux hdparm smartmontools
```

Note: 
- `hdinfo` is _only_ supported on Linux based systems (no Windows or MacOS)
- `hdinfo` has only been tested on a Debian AMD64 system.  If you encounter bugs on another Linux variant or architecture (ARM64), please report it

## Installation

### From source

Requires Go 1.25+.

```bash
go install github.com/cgons/hdinfo/cmd/hdinfo@latest
```

### Build manually

```bash
git clone https://github.com/cgons/hdinfo.git
cd hdinfo
go build -o hdinfo ./cmd/hdinfo
```

## Usage

> **Note:** `hdinfo` requires root privileges. All commands must be run with `sudo`.

### Global Options

| Flag | Description |
|------|-------------|
| `--no-color` | Disable colorized output |

### `hdinfo disks`

List all system disks and their details.

```
sudo hdinfo disks [options]
```

**Options:**

| Flag | Alias | Description |
|------|-------|-------------|
| `--smart-data` | | Include SMART data (temperature, power-on hours, power cycle count) |
| `--force` | | Wake sleeping disks to fetch SMART data (requires `--smart-data`) |
| `--silent` | `-s` | Suppress informational hints in the output |
| `--no-stats` | | Hide disk statistics totals from output |
| `--help` | `-h` | Show help |

**Examples:**

```bash
# List all disks
sudo hdinfo disks

# Include SMART data (skips sleeping disks)
sudo hdinfo disks --smart-data

# Include SMART data and wake sleeping disks
sudo hdinfo disks --smart-data --force

# Silent mode (no hints)
sudo hdinfo disks -s --smart-data

# Disable colors (useful when piping to a file)
sudo hdinfo disks --no-color
```

**Sample output:**

```
Name   Model              Capacity  IsSSD  Interface  State
sda    WDC WD100EMAZ-00   9.1T      false  sata       active/idle
sdb    WDC WD100EMAZ-00   9.1T      false  sata       active/idle
nvme0  Samsung 990 Pro    1.8T      true   nvme       active/idle
```

### `hdinfo mounts`

List disk partitions and their mount points with space usage.

```
sudo hdinfo mounts
```

**Options:**

| Flag | Alias | Description |
|------|-------|-------------|
| `--help` | `-h` | Show help |

**Example:**

```bash
sudo hdinfo mounts
```

**Sample output:**

```
Name    Model              MountPoint   Capacity  UsedSpace  FreeSpace  Used/Free
sda1    WDC WD100EMAZ-00   /mnt/data1   9.1T      4.5T       4.6T       50% / 50%
sdb1    WDC WD100EMAZ-00   /mnt/data2   9.1T      1.2T       7.9T       13% / 87%
nvme0p1 Samsung 990 Pro    /            1.8T      200G       1.6T       11% / 89%
```

## AI Usage
AI usage was limited to the following areas:
- README updates
- Generation of regex patterns to parse disk info
- Updates to colorized output
- Addition of conveniece flags (--no-color, --no-stats)

## License

[MIT](https://github.com/cgons/hdinfo/blob/master/LICENSE)

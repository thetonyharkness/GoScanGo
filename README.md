# GoScanGo
GoScanGo is a simple port scanning tool written in Go that performs TCP port scans on specified targets and generates a report of open ports. It utilizes the net package for basic port scanning and the os/exec package to run Nmap for detailed port analysis.

## Features
- Scan multiple targets for open TCP ports.
- Optionally specify a range or list of ports to scan.
- Run Nmap on open ports for detailed service detection.

## Installation
It is recommended to download pre-compiled binaries for amd64 Linux/Windows and arm64 Darwin from the [latest releases](https://github.com/thetonyharkness/goscango/releases)

### Install From Source
```bash
$ go get github.com/thetonyharkness/goscango
```

## Usage
```text
$ ./goscango -h


  _____       _____                  _____
 / ____|     / ____|                / ____|
| |  __  ___| (___   ___ __ _ _ __ | |  __  ___
| | |_ |/ _ \\___ \ / __/ _\ | '_ \| | |_ |/ _ \
| |__| | (_) |___) | (_| (_| | | | | |__| | (_) |
 \_____|\___/_____/ \___\__,_|_| |_|\_____|\___/



Usage:
$ goscango [OPTIONS] targets

Options:
-A             Additional arguments for Nmap scan
-h             Print usage information
-o <filename>  Output file
-p             Ports to scan (comma-separated or range)
-t             File containing targets
--timeout      Timeout in milliseconds (default: 1500)
-V             Print version information


Examples:
goscango 10.10.10.100
goscango -p 80,443,8080 10.10.10.100
goscango -o results.txt -p 1-1024 10.10.10.100
goscango -A '-sV -sC' 10.10.10.100
```

To use goscango, you can run it from the command line with the following options:

```bash
$ go run goscango.go [OPTIONS] targets
```

## Examples
### Scan all ports for a single target
```bash
$ go run goscango.go 192.168.1.1
```
### Scan specific ports for a single target
```bash
$ go run goscango.go -p 80,443,8080 192.168.1.1
```

### Scan a range of ports for a single target
```bash
$ go run goscango.go -p 1-1024 192.168.1.1
```
### Scan all ports for multiple targets
```bash
$ go run goscan.go -t targets.txt
```

## License
This project is licensed under the GNU General Public License - see the LICENSE file for details.

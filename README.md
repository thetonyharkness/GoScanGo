# goscango
goscango is a simple port scanning tool written in Go that performs TCP port scans on specified targets and generates a report of open ports. It utilizes the net package for basic port scanning and the os/exec package to run Nmap for detailed port analysis.

## Features
- Scan multiple targets for open TCP ports.
- Optionally specify a range or list of ports to scan.
- Run Nmap on open ports for detailed service detection.

## Installation
It is recommended to download pre-compiled binaries for amd64 Linux/Windows and arm64 Darwin from the latest releases

## Usage
To use goscango, you can run it from the command line with the following options:

```bash
$ go run goscango.go [OPTIONS] targets
```
## Options
- -t: File containing targets.
- -o: Output file for scan results.
- --timeout: Timeout in milliseconds (default: 1500).
- -p: Ports to scan (comma-separated or range).
- -V: Print version information.
- -h: Print usage information.

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
This project is licensed under the GNU Public License - see the LICENSE file for details.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	// Define command-line flags
	targetsFile := flag.String("t", "", "File containing targets")
	outputFile := flag.String("o", "", "Output file")
	timeout := flag.Int("timeout", 1500, "Timeout in milliseconds")
	portsArg := flag.String("p", "", "Ports to scan (comma-separated or range)")
	version := flag.Bool("V", false, "Print version information")
	help := flag.Bool("h", false, "Print usage information")

	// Parse command-line flags
	flag.Parse()

	if *targetsFile == "" && len(flag.Args()) == 0 && !*help {
		printUsage()
		return
	}

	if *version {
		fmt.Println("goscango version 1.0")
		return
	}

	if *help {
		printUsage()
		return
	}

	// If targets file is provided, read targets from the file
	var targets []string
	if *targetsFile != "" {
		file, err := os.Open(*targetsFile)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			targets = append(targets, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error:", err)
			return
		}
	} else {
		// Otherwise, use command-line targets
		targets = flag.Args()
	}

	if len(targets) == 0 {
		fmt.Println("Error: No targets specified")
		return
	}

	// Perform port scan for each target
	var wg sync.WaitGroup
	for _, target := range targets {
		wg.Add(1)
		go func(target string) {
			defer wg.Done()
			scanTarget(target, *timeout, *outputFile, *portsArg)
		}(target)
	}
	wg.Wait()

	fmt.Println("Scan complete")
}

func scanTarget(target string, timeout int, outputFile, portsArg string) {
	// Split target into IP and port range
	parts := strings.Split(target, ":")
	ip := parts[0]
	ports := "1-65535"
	if len(parts) > 1 {
		ports = parts[1]
	}

	// If portsArg is provided, use it for port scanning
	if portsArg != "" {
		ports = portsArg
	}

	// Perform port scan
	openPorts := scanPorts(ip, ports, timeout)

	// Run Nmap on open ports
	nmapOutput := runNmap(ip, openPorts)

	// Write results to the specified output file
	if outputFile != "" {
		writeResults(outputFile, ip, openPorts, nmapOutput)
	}
}

func scanPorts(ip, ports string, timeout int) []string {
	fmt.Printf("Scanning ports for %s...\n", ip)
	results := []string{}

	// Split ports into individual ports or range
	portRanges := strings.Split(ports, ",")
	for _, portRange := range portRanges {
		// Check if the portRange is a range (e.g., 80-100)
		if strings.Contains(portRange, "-") {
			// Split range into start and end
			startEnd := strings.Split(portRange, "-")
			start, _ := strconv.Atoi(startEnd[0])
			end, _ := strconv.Atoi(startEnd[1])

			// Perform port scan for each port in the range
			for port := start; port <= end; port++ {
				target := fmt.Sprintf("%s:%d", ip, port)
				if isOpen(target, time.Duration(timeout)*time.Millisecond) {
					fmt.Printf("Open %s\n", target)
					results = append(results, strconv.Itoa(port))
				}
			}
		} else {
			// Single port scan
			port, _ := strconv.Atoi(portRange)
			target := fmt.Sprintf("%s:%d", ip, port)
			if isOpen(target, time.Duration(timeout)*time.Millisecond) {
				fmt.Printf("Open %s\n", target)
				results = append(results, portRange)
			}
		}
	}

	return results
}

func isOpen(target string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func runNmap(ip string, openPorts []string) string {
	if len(openPorts) == 0 {
		fmt.Printf("No open ports found for %s\n", ip)
		return ""
	}

	fmt.Printf("Running Nmap for %s on ports %s\n", ip, strings.Join(openPorts, ","))

	// Run Nmap command
	cmd := exec.Command("nmap", "-Pn", "-vvv", "-p", strings.Join(openPorts, ","), ip)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running Nmap command:", err)
		return ""
	}

	// Print Nmap output
	fmt.Println(string(stdout))
	return string(stdout)
}

func writeResults(outputFile, ip string, openPorts []string, nmapOutput string) {
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening output file:", err)
		return
	}
	defer file.Close()

	// Write open ports to file
	for _, port := range openPorts {
		result := fmt.Sprintf("Open %s:%s\n", ip, port)
		_, err := file.WriteString(result)
		if err != nil {
			fmt.Println("Error writing to output file:", err)
			return
		}
	}

	// Write Nmap output to file
	_, err = file.WriteString("\n" + nmapOutput)
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		return
	}
}

func printUsage() {
	fmt.Println(`
Usage:
goscango [OPTIONS] targets

Options:
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
`)
}

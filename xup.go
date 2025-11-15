package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

type NmapRun struct {
	XMLName xml.Name `xml:"nmaprun"`
	Hosts   []Host   `xml:"host"`
}

type Host struct {
	XMLName xml.Name `xml:"host"`
	Status  Status   `xml:"status"`
	Address Address  `xml:"address"`
	Ports   Ports    `xml:"ports"`
}

type Status struct {
	State string `xml:"state,attr"`
}

type Address struct {
	Addr     string `xml:"addr,attr"`
	AddrType string `xml:"addrtype,attr"`
}

type Ports struct {
	Port []Port `xml:"port"`
}

type Port struct {
	Protocol string `xml:"protocol,attr"`
	PortID   string `xml:"portid,attr"`
	State    State  `xml:"state"`
}

type State struct {
	State string `xml:"state,attr"`
}

func main() {
	help := flag.Bool("help", false, "Show help")
	onlyIP := flag.Bool("onlyip", false, "Print only IP addresses")
	flag.Parse()

	args := flag.Args()

	// Check if -onlyip appears in args (when used after mode)
	mode := ""
	for i, arg := range args {
		if arg == "-onlyip" || arg == "--onlyip" {
			*onlyIP = true
			// Remove -onlyip from args
			args = append(args[:i], args[i+1:]...)
			break
		}
		if arg == "masscan" || arg == "nmap" {
			mode = arg
		}
	}

	if mode == "" && len(args) > 0 {
		mode = args[0]
	}

	if mode == "" {
		if *help {
			printHelp()
			return
		}
		fmt.Fprintln(os.Stderr, "Error: Please specify 'masscan' or 'nmap' mode")
		printHelp()
		os.Exit(1)
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error checking stdin:", err)
		return
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		if *help {
			printHelp()
			return
		}
		fmt.Fprintln(os.Stderr, "Error: No input from stdin")
		printHelp()
		os.Exit(1)
	}

	switch mode {
	case "masscan":
		parseMasscan(os.Stdin, *onlyIP)
	case "nmap":
		parseNmap(os.Stdin, *onlyIP)
	default:
		fmt.Fprintln(os.Stderr, "Error: Unknown mode. Use 'masscan' or 'nmap'")
		printHelp()
		os.Exit(1)
	}
}

func parseMasscan(reader io.Reader, onlyIP bool) {
	scanner := bufio.NewScanner(reader)
	re := regexp.MustCompile(`Discovered open port (\d+)/(tcp|udp) on (\d+\.\d+\.\d+\.\d+)`)
	seenIPs := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if match != nil {
			ip := match[3]
			port := match[1]
			if onlyIP {
				if !seenIPs[ip] {
					fmt.Println(ip)
					seenIPs[ip] = true
				}
			} else {
				fmt.Printf("%s:%s\n", ip, port)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
	}
}

func parseNmap(reader io.Reader, onlyIP bool) {
	var nmapRun NmapRun
	decoder := xml.NewDecoder(reader)

	err := decoder.Decode(&nmapRun)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing XML:", err)
		return
	}

	seenIPs := make(map[string]bool)

	for _, host := range nmapRun.Hosts {
		if host.Status.State != "up" {
			continue
		}

		ip := host.Address.Addr
		if ip == "" {
			continue
		}

		if onlyIP {
			if !seenIPs[ip] {
				fmt.Println(ip)
				seenIPs[ip] = true
			}
		} else {
			for _, port := range host.Ports.Port {
				if port.State.State == "open" {
					fmt.Printf("%s:%s\n", ip, port.PortID)
				}
			}
		}
	}
}

func printHelp() {
	fmt.Println("\nUsage:")
	fmt.Println("  cat scope.txt | xup masscan")
	fmt.Println("  cat scope.txt | xup nmap")
	fmt.Println("  cat scope.txt | xup masscan -onlyip")
	fmt.Println("  cat scope.txt | xup nmap -onlyip")
	fmt.Println("\nOptions:")
	fmt.Println("  -help     Show help")
	fmt.Println("  -onlyip   Print only IP addresses")
	fmt.Println("\n@musana | musana.net")
}

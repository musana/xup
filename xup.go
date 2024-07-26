package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

func main() {

	help := flag.Bool("help", false, "Show help")
	onlyIP := flag.Bool("onlyip", false, "Print only IP addresses")
	flag.Parse()

	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error checking stdin:", err)
		return
	}

	if (stat.Mode()&os.ModeCharDevice) != 0 || *help {
		printHelp()
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	re := regexp.MustCompile(`Discovered open port (\d+)/(tcp|udp) on (\d+\.\d+\.\d+\.\d+)`)

	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if match != nil {
			ip := match[3]
			port := match[1]
			if *onlyIP {
				fmt.Println(ip)
			} else {
				fmt.Printf("%s:%s\n", ip, port)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
	}
}

func printHelp() {
	fmt.Println("\nUsage: masscan x.x.x.x/x -p445 | xup")
	fmt.Println("Options:")
	fmt.Println("  -help     Show help")
	fmt.Println("  -onlyip   Print only IP addresses")
	fmt.Println("\n@musana | musana.net")
}

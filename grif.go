package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	toast "gopkg.in/toast.v1"
)

func main() {
	hosts, err := getHosts()
	if err != nil {
		log.Fatal(err)
	}
	for _, host := range hosts {
		notify()
		fmt.Println(host)
		if strings.Contains(host, "http://") || strings.Contains(host, "https://") {
			fmt.Println("Detected HTTP/S")
		}
		if strings.Contains(host, "icmp://") {
			fmt.Println("Detected ICMP")
		}
	}
	notify()
}

func getHosts() ([]string, error) {
	// getHosts reads a 'hosts' file into memory
	// and returns a slice of its hosts
	file, err := os.Open("hosts")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func notify() {
	notification := toast.Notification{
		AppID:   "Grif",
		Title:   "Grif has detected an outage",
		Message: "Host icmp://8.8.8.8 didn't respond after 3 retries",
		Icon:    "assets/grif.ico", // This file must exist (remove this line if it doesn't)
		// Actions: []toast.Action{
		// 	{"protocol", "Acknowledge", ""},
		// 	{"protocol", "Ignore", ""},
		// },
	}
	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}

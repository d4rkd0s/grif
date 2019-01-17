package main

import (
	"bufio"
	"fmt"
	"github.com/gen2brain/beeep"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	hosts, err := getHosts()
	if err != nil {
		log.Fatal(err)
	}
	for _, host := range hosts {
		fmt.Println(host)
		if strings.Contains(host, "http://") || strings.Contains(host, "https://") {
			fmt.Println("Detected HTTP/S")
			checkHTTP(host)
		}
		if strings.Contains(host, "icmp://") {
			fmt.Println("Detected ICMP")
			checkICMP(host)
		}
	}
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

func notify(message string) {
	err := beeep.Notify("Title", message, "assets/grif.png")
	if err != nil {
		panic(err)
	}
}

func checkHTTP(host string) {
	resp, err := http.Get(host)
	if err != nil {
		notify(fmt.Sprintf("%s", err))
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		fmt.Println(fmt.Sprintf("%s replied with a 2xx", host))
	} else {
		notify(fmt.Sprintf("%s returned a non 2xx HTTP response", host))
	}
}

func checkICMP(host string) {

}


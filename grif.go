package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/sparrc/go-ping"
)

// Length in seconds between each cycle of checks
var wait time.Duration = 3
var checking = false
var ready = false
var numOfChecksRan uint32 = 0

func main() {
	fmt.Println("Grif v0.1 - https://github.com/d4rkd0s/grif (c) 2019 d4rkd0s")
	if !checkForAdmin() {
		os.Exit(126)
	} else {
		fmt.Println("Administrator command prompt achieved")
	}
	checkAndCreateHostsFile()
	// Code to ask amount of delay
	//
	//fmt.Print("Enter number of seconds to wait between checking all hosts: ")
	//var i int
	//_, err := fmt.Scanf("%d", &i)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//wait = time.Duration(i)
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon("assets/grif.ico"))
	systray.SetTitle("Grif")
	systray.SetTooltip("Grif is running")
	mQuit := systray.AddMenuItem("Quit", "Quits Grif")
	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()

	go func() {
		for range time.Tick(wait * time.Second) {
			if !ready {
				go notify("Grif is now running")
				bark()
			}
			ready = true
			if !checking {
				go checkHosts()
			} else {
				fmt.Println("Busy checking hosts, waiting another cycle")
			}
		}
	}()
}

func onExit() {
	// Cleaning stuff here.
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}

func checkHosts() {
	checking = true
	fmt.Println("Checking hosts...")
	hosts, err := getHosts()
	fmt.Println(hosts)
	if err != nil {
		log.Fatal(err)
	}
	for _, host := range hosts {
		numOfChecksRan++
		if strings.Contains(host, "http://") || strings.Contains(host, "https://") {
			checkHTTP(host)
		}
		if strings.Contains(host, "icmp://") {
			fmt.Println("ICMP unsupported currently")
			//  checkICMP(host)
		}
	}
	systray.SetTooltip(fmt.Sprintf("%s checks ran, %d hosts up, %d hosts down", strconv.Itoa(int(numOfChecksRan)), getCountOfUpHosts(), getCountOfDownHosts()))
	fmt.Println("Done checking hosts")
	checking = false
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

func getCountOfUpHosts() uint32 {
	var numOfUpHosts uint32 = 0
	var hosts, err = getHosts()
	if err != nil {
		log.Fatal(err)
	}
	for _, host := range hosts {
		if !strings.HasPrefix(host, "#") {
			numOfUpHosts++
		}
	}
	return numOfUpHosts
}

func getCountOfDownHosts() uint32 {
	var numOfDownHosts uint32 = 0
	var hosts, err = getHosts()
	if err != nil {
		log.Fatal(err)
	}
	for _, host := range hosts {
		if strings.HasPrefix(host, "#") {
			numOfDownHosts++
		}
	}
	return numOfDownHosts
}

func checkIfHostWasUp(host string) bool {
	hosts, err := getHosts()
	if err != nil {
		log.Fatal(err)
	}
	for _, hostFromFile := range hosts {
		if strings.Contains(hostFromFile, host) {
			if strings.HasPrefix(hostFromFile, "#") {
				return false
			} else {
				return true
			}
		}
	}
	return false
}

func markHostDown(host string) {
	hosts, err := getHosts()
	if err != nil {
		log.Fatal(err)
	}
	for i, hostFromFile := range hosts {
		if strings.Contains(hostFromFile, host) {
			if !strings.HasPrefix(hostFromFile, "#") {
				hosts[i] = fmt.Sprintf("#%s", hostFromFile)
				output := strings.Join(hosts, "\n")
				err = ioutil.WriteFile("hosts", []byte(output), 0644)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}
	return
}

func markHostUp(host string) {
	hosts, err := getHosts()
	if err != nil {
		log.Fatal(err)
	}
	for i, hostFromFile := range hosts {
		if strings.Contains(hostFromFile, host) {
			if strings.HasPrefix(hostFromFile, "#") {
				hosts[i] = fmt.Sprintf("%s", host)
				output := strings.Join(hosts, "\n")
				err = ioutil.WriteFile("hosts", []byte(output), 0644)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}
	return
}

func notify(message string) {
	// If the host is commented out, don't alert
	err := beeep.Notify("Grif", message, "assets/grif.png")
	if err != nil {
		panic(err)
	}
	// Wait 2 seconds before continuing to let the notifications pool up slowly
	time.Sleep(2 * time.Second)
	return
}

func alert(message string, host string) {
	// If the host is commented out, don't alert
	if checkIfHostWasUp(host) {
		err := beeep.Alert("Grif", message, "assets/grif.png")
		if err != nil {
			panic(err)
		}
		bark()
		markHostDown(host)
		// TODO: Add downed hosts to list in right click menu
		// If a user clicked a downed host, Grif will recheck (ie remove #)
	}
	// Wait 2 seconds before continuing to let the notifications pool up slowly
	time.Sleep(2 * time.Second)
	return
}

func checkHTTP(host string) {
	if strings.HasPrefix(host, "#") {
		host = trimFirstRune(host)
	}

	//resp, err := http.Get(host)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(host)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err), host)
		alert(fmt.Sprintf("%s", err), host)
	} else {
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			fmt.Println(fmt.Sprintf("%s replied with a 2xx", host))
			markHostUp(host)
		} else {
			alert(fmt.Sprintf("%s returned a non 2xx HTTP response", host), host)
		}
	}
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func checkICMP(host string) {
	if strings.HasPrefix(host, "#") {
		host = trimFirstRune(host)
	}
	host = strings.TrimPrefix(host, "icmp://")
	pinger, err := ping.NewPinger(host)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err), host)
		alert(fmt.Sprintf("%s", err), host)
	}
	pinger.Count = 1
	pinger.Run() // blocks until finished

	//if err != nil {
	//
	//} else {
	//	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
	//		fmt.Println(fmt.Sprintf("%s replied with a 2xx", host))
	//		markHostUp(host)
	//	} else {
	//		alert(fmt.Sprintf("%s returned a non 2xx HTTP response", host), host)
	//	}
	//}
}

func bark() {
	// Open first sample File
	f, err := os.Open("assets/bark.mp3")

	// Check for errors when opening the file
	if err != nil {
		log.Fatal(err)
	}

	// Decode the .mp3 File, if you have a .wav file, use wav.Decode(f)
	s, format, _ := mp3.Decode(f)

	// Init the Speaker with the SampleRate of the format and a buffer size of 1/10s
	err2 := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err2 != nil {
		panic(err2)
	}

	// Channel, which will signal the end of the playback.
	playing := make(chan struct{})

	// Grif volume to half
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	volume.Volume += 0.5

	// Now we Play our Streamer on the Speaker
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		// Callback after the stream Ends
		close(playing)
	})))
	<-playing
}

func checkAndCreateHostsFile() bool {
	_, err := os.Open("hosts")
	if err != nil {
		_, err2 := os.OpenFile("hosts", os.O_RDONLY|os.O_CREATE, 0666)
		if err2 != nil {
			fmt.Println("Unable to create hosts file in current directory, we have admin rights, so something must be wrong.")
			log.Fatal(err2)
		} else {
			fmt.Println("Created hosts file in current directory")
		}
	} else {
		fmt.Println("Detected hosts file in current directory")
	}
	return true
}

func checkForAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}


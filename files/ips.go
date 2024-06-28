package files

import (
	"bufio"
	"fmt"
	"os"
)

func ReadIPsFile(filename string) []string {
	// read a file and by every new line return a string and append it to the ips array

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Failed to open", filename)
		fmt.Println(err)
		return []string{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	ips := []string{}

	for scanner.Scan() {
		ips = append(ips, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return []string{}
	}

	return ips
}

func SaveStringToFile(content string) error {
	file, err := os.OpenFile("vuln.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, content)
	if err != nil {
		return err
	}

	return nil
}

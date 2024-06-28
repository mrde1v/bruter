package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func ClearTerminal() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	case "linux":
		cmd = exec.Command("clear")
	default:
		fmt.Println("Your OS is not supported")
		os.Exit(1)
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

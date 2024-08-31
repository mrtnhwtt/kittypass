package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.design/x/clipboard"
)

func ParseTimestamp(timestamp string) (string, error) {
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return "", fmt.Errorf("error parsing time: %s", err)
	}
	return parsedTime.Format("02 Jan 2006 15:04"), nil
}

// IsWSL checks if the current environment is WSL
func IsWSL() bool {
	file, err := os.Open("/proc/sys/kernel/osrelease")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		release := scanner.Text()
		if strings.Contains(release, "Microsoft") || strings.Contains(release, "WSL") {
			return true
		}
	}
	return false
}

func AddToClipboard(text string) error {
	if IsWSL() {
		// If in WSL, use clipboard.exe because x/clipboard panics or doesn't work on WSL
		cmd := exec.Command("clip.exe")
		cmd.Stdin = strings.NewReader(text)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to copy to clipboard using clip.exe: %v", err)
		}
		return nil
	}
	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("failed to initialize clipboard: %v", err)
	}
	clipboard.Write(clipboard.FmtText, []byte(text))

	return nil
}

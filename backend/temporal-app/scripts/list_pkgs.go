package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("go", "list", "./...")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	lines := bytes.Split(out, []byte("\n"))
	for _, line := range lines {
		s := string(line)
		if s != "" && !strings.Contains(s, "/tests/mocks") {
			fmt.Print(s + " ")
		}
	}
}

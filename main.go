package main

import (
	"os"

	"github.com/br-lemes/esol/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

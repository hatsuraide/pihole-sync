package main

import (
	"os"

	"github.com/Kizuren/pihole-sync/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

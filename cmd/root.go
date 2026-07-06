package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Kizuren/pihole-sync/internal/log"
	"github.com/Kizuren/pihole-sync/version"
)

var rootCmd = &cobra.Command{
	Use:     "pihole-sync",
	Version: version.Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(log.Init)
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

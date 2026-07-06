package cmd

import (
	kservice "github.com/kardianos/service"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/Kizuren/pihole-sync/internal/daemon"
)

var serviceEnvFile string

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Install and control pihole-sync as a native background service (no Docker required)",
}

var serviceInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install pihole-sync as a service that starts on boot and restarts automatically",
	Run: func(cmd *cobra.Command, args []string) {
		runServiceAction(func(s kservice.Service) error { return s.Install() }, "installed")
	},
}

var serviceUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove the installed pihole-sync service",
	Run: func(cmd *cobra.Command, args []string) {
		runServiceAction(func(s kservice.Service) error { return s.Uninstall() }, "uninstalled")
	},
}

var serviceStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the installed pihole-sync service",
	Run: func(cmd *cobra.Command, args []string) {
		runServiceAction(func(s kservice.Service) error { return s.Start() }, "started")
	},
}

var serviceStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running pihole-sync service",
	Run: func(cmd *cobra.Command, args []string) {
		runServiceAction(func(s kservice.Service) error { return s.Stop() }, "stopped")
	},
}

var serviceRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the pihole-sync service",
	Run: func(cmd *cobra.Command, args []string) {
		runServiceAction(func(s kservice.Service) error { return s.Restart() }, "restarted")
	},
}

var serviceStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the pihole-sync service status",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := daemon.New(serviceEnvFile)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to configure service")
		}

		status, err := s.Status()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get service status")
		}

		log.Info().Msg(statusLabel(status))
	},
}

// serviceRunCmd is the entrypoint the OS service manager (systemd/launchd/Windows SCM) actually
// executes once installed; it is not meant to be invoked directly by users, who should use `run`
// for foreground/Docker use or `service install`/`service start` for a native background service.
var serviceRunCmd = &cobra.Command{
	Use:    "service-run",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := daemon.New(serviceEnvFile)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to configure service")
		}

		if err := s.Run(); err != nil {
			log.Fatal().Err(err).Msg("Service exited with error")
		}
	},
}

func runServiceAction(action func(s kservice.Service) error, verb string) {
	s, err := daemon.New(serviceEnvFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to configure service")
	}

	if err := action(s); err != nil {
		log.Fatal().Err(err).Msgf("Failed to run service action")
	}

	log.Info().Msgf("Service %s", verb)
}

func statusLabel(status kservice.Status) string {
	switch status {
	case kservice.StatusRunning:
		return "running"
	case kservice.StatusStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(serviceRunCmd)
	serviceCmd.AddCommand(serviceInstallCmd, serviceUninstallCmd, serviceStartCmd, serviceStopCmd, serviceRestartCmd, serviceStatusCmd)

	serviceCmd.PersistentFlags().StringVar(&serviceEnvFile, "env-file", "",
		"Path to a .env file the service should load on start (recommended over baking secrets into the service definition)")
	serviceRunCmd.Flags().StringVar(&serviceEnvFile, "env-file", "", "Read env from `.env` file")
}

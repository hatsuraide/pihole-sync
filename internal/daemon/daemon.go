// Package daemon lets pihole-sync run as a native background service (systemd, launchd or
// Windows Service) instead of depending on Docker/an external scheduler for continuous operation.
package daemon

import (
	kservice "github.com/kardianos/service"
	"github.com/rs/zerolog/log"

	"github.com/Kizuren/pihole-sync/internal/config"
	svc "github.com/Kizuren/pihole-sync/internal/service"
)

const Name = "pihole-sync"

// RunCommand is the hidden CLI command the OS service manager invokes to actually run the sync
// service; it is what daemon.New registers as the service's Arguments.
const RunCommand = "service-run"

// New builds the service handle used both to control the OS service (install/start/stop/status)
// and, via serviceRun, as the payload the service manager actually executes.
func New(envFile string) (kservice.Service, error) {
	svcConfig := &kservice.Config{
		Name:        Name,
		DisplayName: "Pi-hole Sync",
		Description: "Synchronizes Pi-hole v6 configuration to replicas on a schedule.",
		Arguments:   runArguments(envFile),
	}

	return kservice.New(&program{envFile: envFile}, svcConfig)
}

func runArguments(envFile string) []string {
	args := []string{RunCommand}
	if envFile != "" {
		args = append(args, "--env-file", envFile)
	}
	return args
}

type program struct {
	envFile string
}

func (p *program) Start(_ kservice.Service) error {
	go p.run()
	return nil
}

func (*program) Stop(_ kservice.Service) error {
	return nil
}

func (p *program) run() {
	if p.envFile != "" {
		if err := config.LoadEnvFile(p.envFile); err != nil {
			log.Fatal().Err(err).Msg("Failed to load env file")
		}
	}

	service, err := svc.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize service")
	}

	if err = service.Run(); err != nil {
		log.Fatal().Err(err).Msg("Sync failed")
	}
}

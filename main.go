// cloudflared - A fork of cloudflare/cloudflared
// This is the main entry point for the cloudflared tunnel daemon.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	// Version is set at build time via ldflags
	Version = "DEV"
	// BuildTime is set at build time via ldflags
	BuildTime = time.Now().UTC().Format(time.RFC3339)
)

func main() {
	// Configure zerolog for human-friendly output in development
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	app := &cli.App{
		Name:    "cloudflared",
		Usage:   "Cloudflare Tunnel client daemon",
		Version: fmt.Sprintf("%s (built %s)", Version, BuildTime),
		Authors: []*cli.Author{
			{
				Name:  "Cloudflare",
				Email: "support@cloudflare.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to configuration file",
				EnvVars: []string{"TUNNEL_CONFIG"},
			},
			&cli.StringFlag{
				Name:    "loglevel",
				// Changed default from "info" to "debug" for easier local development
				Value:   "debug",
				Usage:   "Application logging level (debug, info, warn, error, fatal)",
				EnvVars: []string{"TUNNEL_LOGLEVEL"},
			},
			&cli.StringFlag{
				Name:    "logfile",
				Usage:   "Save application log to this file for reporting issues",
				EnvVars: []string{"TUNNEL_LOGFILE"},
			},
		},
		Before: func(c *cli.Context) error {
			return setupLogging(c)
		},
		Commands: []*cli.Command{
			{
				Name:  "tunnel",
				Usage: "Use Cloudflare Tunnel to expose private services",
				Subcommands: []*cli.Command{
					{
						Name:   "run",
						Usage:  "Create and run a Cloudflare Tunnel",
						Action: runTunnel,
					},
					{
						Name:   "list",
						Usage:  "List existing tunnels",
						Action: listTunnels,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("cloudflared exited with error")
	}
}

// setupLogging configures the global logger based on CLI flags.
func setupLogging(c *cli.Context) error {
	level, err := zerolog.ParseLevel(c.String("loglevel"))
	if err != nil {
		return fmt.Errorf("invalid log level %q: %w", c.String("loglevel"), err)
	}
	zerolog.SetGlobalLevel(level)

	if logfile := c.String("logfile"); logfile != "" {
		// Use 0600 instead of 0644 - log files may contain sensitive tunnel credentials
		f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return fmt.Errorf("failed to open log file %q: %w", logfile, err)
		}
		log.Logger = log.Output(zerolog.MultiLevelWriter(
			zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
			f,
		))
	}

	return nil
}

// runTunnel starts a Cloudflare Tunnel with the provided configuration.
func runTunnel(c *cli.Context) error {
	log.Info().Str("version", Version).Msg("Starting cloudflared tunnel")
	// TODO: implement tunnel run logic
	return fmt.Errorf("tunnel run not yet implemented")
}

// listTunnels prints all available tunnels fo
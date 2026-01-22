package main

import (
	"log/slog"
	"os"

	"github.com/oxiginedev/sabipass/cmd/http"
	"github.com/oxiginedev/sabipass/config"
	"github.com/spf13/cobra"
)

func main() {
	err := os.Setenv("TZ", "")
	if err != nil {
		slog.Error("cmd: failed to set timezone", slog.Any("error", err))
		os.Exit(1)
	}

	cfg := &config.Config{}

	rootCmd := &cobra.Command{
		Use:   "sabipass",
		Short: "Fun realtime games with friends and family",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			envFile, err := cmd.Flags().GetString("env")
			if err != nil {
				return err
			}

			return config.Load(envFile, cfg)
		},
	}

	rootCmd.PersistentFlags().String("env", ".env", "Environment file")

	rootCmd.AddCommand(http.Command(cfg))

	err = rootCmd.Execute()
	if err != nil {
		slog.Error("cmd: failed to execute root command", slog.Any("error", err))
		os.Exit(1)
	}
}

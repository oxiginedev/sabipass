package http

import (
	"log/slog"
	"os"

	"github.com/oxiginedev/sabipass/config"
	"github.com/oxiginedev/sabipass/internal/api"
	"github.com/oxiginedev/sabipass/internal/database/postgres"
	"github.com/oxiginedev/sabipass/internal/pkg/jwt"
	"github.com/oxiginedev/sabipass/internal/server"
	"github.com/spf13/cobra"
)

func Command(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Start the HTTP server",
		Run: func(cmd *cobra.Command, args []string) {
			pgdb, err := postgres.NewDB(cfg)
			if err != nil {
				slog.Error("could not connect to database", slog.Any("error", err))
				os.Exit(1)
			}

			userRepo := postgres.NewUserRepository(pgdb)
			quizRepo := postgres.NewQuizRepository(pgdb)
			questionTypeRepo := postgres.NewQuestionTypeRepository(pgdb)

			tokenManager := jwt.NewJwtTokenManager(cfg)
			handler := api.NewAPI(cfg, tokenManager, userRepo, quizRepo, questionTypeRepo)

			srv := server.NewServer(cfg, func() {
				err := pgdb.Close()
				if err != nil {
					slog.Error("could not close database connection", slog.Any("error", err))
				}
			})
			srv.SetHandler(handler.RegisterRoutes())
			srv.Listen()
		},
	}
}

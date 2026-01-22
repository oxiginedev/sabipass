package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/oxiginedev/sidekik"
)

// ENUM(production, local)
type Environment string

type Config struct {
	Environment Environment
	HTTP        struct {
		Port uint16 `default:"8000"`
	}

	Database struct {
		Postgres struct {
			DSN          string        `envconfig:"SABIPASS_POSTGRES_DSN"`
			QueryTimeout time.Duration `envconfig:"SABIPASS_POSTGRES_QUERY_TIMEOUT" default:"5s"`
		}
	}

	Oauth struct {
		Google struct {
			ClientID     string `envconfig:"SABIPASS_GOOGLE_CLIENT_ID"`
			ClientSecret string `envconfig:"SABIPASS_GOOGLE_CLIENT_SECRET"`
			RedirectURL  string `envconfig:"SABIPASS_GOOGLE_REDIRECT_URL"`
		}
	}

	Auth struct {
		JWT struct {
			SecretKey string
			Expiry    time.Duration `default:"1h"`
		}
	}
}

func Load(pathToFile string, cfg *Config) error {
	var err error
	if !sidekik.IsStringEmpty(pathToFile) {
		err = godotenv.Load(pathToFile)
	} else {
		err = godotenv.Load()
	}

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return envconfig.Process("SABIPASS", cfg)
}

package cli

import (
	"context"
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/zlobste/binance-mystery-buyer/internal/config"
	"github.com/zlobste/binance-mystery-buyer/internal/services/buyer"
	"os"
)

func Run(args []string) bool {
	cfg := config.New(os.Getenv("CONFIG"))
	ctx := context.Background()
	log := cfg.Log()

	defer func() {
		if rvr := recover(); rvr != nil {
			log.Error(fmt.Sprintf("app panicked: %s", rvr))
		}
	}()

	app := kingpin.New("binance-mystery-buyer", "")

	runCmd := app.Command("run", "run command")
	buyerCmd := runCmd.Command("buyer", "run a service to deposit into Odin")

	cmd, err := app.Parse(args[1:])
	if err != nil {
		log.WithError(err).Error("failed to parse arguments")

		return false
	}

	switch cmd {
	case buyerCmd.FullCommand():
		if err := buyer.New(cfg).Run(ctx); err != nil {
			log.WithError(err).Error("failed to run deploy service")

			return false
		}
	default:
		log.WithField("command", cmd).Error("unknown command")

		return false
	}

	return true
}

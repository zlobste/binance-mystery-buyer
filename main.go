package main

import (
	"github.com/zlobste/binance-mystery-buyer/internal/cli"
	"os"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}

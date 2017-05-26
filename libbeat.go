package main

import (
	"os"

	"github.com/dearcode/libbeat/beat"
	"github.com/dearcode/libbeat/mock"
)

func main() {
	if err := beat.Run(mock.Name, mock.Version, mock.New); err != nil {
		os.Exit(1)
	}
}

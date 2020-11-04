package main

import (
	"github.com/refs/panic-middleware/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	golog "log"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	s := service.Service{}
	golog.Fatal(s.Run())
}

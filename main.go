package main

import (
	"github.com/refs/panic-middleware/service"
	"log"
)

func main() {
	s := service.Service{}
	log.Fatal(s.Run())
}

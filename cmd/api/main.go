package main

import (
	"github.com/scalarorg/scalar-healer/cmd/api/server"
	"github.com/scalarorg/scalar-healer/config"
)

func main() {
	config.LoadEnv()
	s := server.New()
	err := s.Start()
	defer s.Close()
	if err != nil {
		panic(err)
	}
}

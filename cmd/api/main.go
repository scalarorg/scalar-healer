package main

import (
	"github.com/scalarorg/scalar-healer/cmd/api/server"
)

func main() {
	s := server.New()
	err := s.Start()
	defer s.Close()
	if err != nil {
		panic(err)
	}
}

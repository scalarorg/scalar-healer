package main

import (
	"github.com/0xdavid7/goes-template/cmd/api/server"
)

func main() {
	s := server.New()
	err := s.Start()
	defer s.Close()
	if err != nil {
		panic(err)
	}
}

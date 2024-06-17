package main

import (
	"log"

	"github.com/msarifin29/be_budget_in/internal/config"
	"github.com/msarifin29/be_budget_in/internal/delivery/controller"
)

var env, port string

func init() {
	env = "dev"
	port = ":8080"
}

func main() {
	logger := config.NewLogger()
	con, errCon := config.LoadConfig("/home/samsul-dev/projects/budget_in_backend", env)
	if errCon != nil {
		log.Fatalf("cannot load config %e :", errCon)
	}
	server, err := controller.NewServer(logger, con, env)
	if err != nil {
		log.Fatalf("cannot create server %e :", errCon)
	}
	errStart := server.Start(con.ServerAddress + port)
	if errStart != nil {
		log.Fatalf("cannot start server %e :", errStart)
	}
}

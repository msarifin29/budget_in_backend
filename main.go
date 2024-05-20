package main

import (
	"log"

	"github.com/msarifin29/be_budget_in/internal/config"
	"github.com/msarifin29/be_budget_in/internal/delivery/controller"
)

func main() {
	logger := config.NewLogger()
	con, errCon := config.LoadConfig(".", "prod")
	if errCon != nil {
		log.Fatalf("cannot load config %e :", errCon)
	}
	server, err := controller.NewServer(logger, con)
	if err != nil {
		log.Fatalf("cannot create server %e :", errCon)
	}
	errStart := server.Start(con.ServerAddress)
	if errStart != nil {
		log.Fatalf("cannot start server %e :", errStart)
	}
}

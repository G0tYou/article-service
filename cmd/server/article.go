package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"article/config"
	"article/pkg/adding"
	"article/pkg/delivery/http/rest"
	"article/pkg/listing"
	"article/pkg/repository/mysql"
	"article/pkg/repository/redis"
)

func main() {
	goEnv := strings.ToLower(os.Getenv("GO_ENV"))
	if goEnv == "" {
		goEnv = "local"
	}

	// Load config
	err := config.Load("config/config.local.yaml")
	if err != nil {
		log.Fatal("Error: Config failed to load - ", err)
	}

	// Run the server
	run(goEnv)
}

func run(goEnv string) {
	// MySQL setup
	mysql, err := mysql.NewStorage(config.My)
	if err != nil {
		log.Fatal("Error: Database failed to connect (", config.My.DSN, ") - ", err)
	}

	// MySQL setup
	redis, err := redis.NewStorage(config.Rd)
	if err != nil {
		log.Fatal("Error: Redis failed to connect (", config.Rd.Addr, ") - ", err)
	}

	// Handler setup
	adder := adding.NewService(mysql, redis)
	lister := listing.NewService(mysql, redis)

	r := rest.Handler(adder, lister)

	host := config.Glb.Serv.Host

	log.Println("Server Running on", goEnv, "environment, listening on", host+":"+config.Serv.Port)
	log.Fatal("Error: Server failed to run - ", http.ListenAndServe(host+":"+config.Serv.Port, r))
}

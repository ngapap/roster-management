package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/pressly/goose"
	"github.com/spf13/viper"

	"roster-management/pkg/postgres"
)

func main() {
	configPath := flag.String("config", "configs/local.env", "config file")
	migrationPath := flag.String("migrations", "migrations/sql", "migrations path")

	config := viper.New()
	config.AutomaticEnv()
	config.SetConfigFile(*configPath)
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	db, err := postgres.Open(config)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("err when closing db connection: %v", err)
		}
	}(db)
	if err != nil {
		log.Panic(err.Error())
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Panic(err.Error())
	}

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		log.Fatal("expected migrations arg")
	}

	goose.SetSequential(true)
	if err = goose.Run(args[0], db, *migrationPath, args[1:]...); err != nil {
		log.Panic(err.Error())
	}
}

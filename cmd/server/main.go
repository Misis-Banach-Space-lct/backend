package main

import (
	"fmt"
	"lct/internal/api/router"
	"lct/internal/config"
	"lct/internal/database"
	"lct/internal/logging"
	"os"
)

// TODO: add graceful shutdown
func main() {
	r, err := setup()
	if err != nil {
		fmt.Printf("error while setting application up: %+v", err)
		os.Exit(-1)
	}

	r.Run()
}

func setup() (*router.Router, error) {
	if err := config.NewConfig(); err != nil {
		return nil, err
	}

	if err := logging.NewLogger(); err != nil {
		return nil, err
	}

	if err := database.NewPostgres(10); err != nil {
		return nil, err
	}

	r, err := router.NewRouter()
	if err != nil {
		return nil, err
	}

	return r, nil
}

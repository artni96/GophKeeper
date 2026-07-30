package main

import (
	"context"
	"log"

	"github.com/artni96/GophKeeper/internal/config"
)

func main() {
	ctx := context.Background()
	cfg, err := config.ParseConfig()

	if err != nil {
		log.Fatal(err)
	}
	err = run(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"context"
	"log"
	"time"

	"github.com/edkuzhakhmetov/go_final_project/internal/application"
)

func main() {
	parent := context.Background()

	ctx, cancel := context.WithTimeout(parent, time.Second*10)
	defer cancel()

	app := application.New()
	err := app.Run(ctx)
	if err != nil {
		log.Fatalln("failed to start app")
	}

}

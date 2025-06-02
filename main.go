package main

import (
	"context"
	"log"

	"github.com/edkuzhakhmetov/go_final_project/internal/application"
)

func main() {
	ctx := context.TODO()
	app := application.New()
	err := app.Run(ctx)
	if err != nil {
		log.Fatalln("failed to start app")
	}

}

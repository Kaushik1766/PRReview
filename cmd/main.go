package main

import "Kaushik1766/PRReview/internal/app"

func main() {
	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	app.Run()
}

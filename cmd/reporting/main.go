package main

import "powerbi-live-reporting/internal/app"

func main() {
	application := app.CreateApp()
	application.Run()
}
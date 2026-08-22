package main

import "github.com/vandad1901/p3s/apps/api/internal/app"

func main() {
	a := app.Boot()
	a.Run()
}

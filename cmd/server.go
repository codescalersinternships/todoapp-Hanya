package main

import (
	"log"
	todo "todo/internal"
)

func main() {
	dbPath, port, err := todo.ParseCommandLine()
	if err != nil {
		log.Fatal(err)
	}
	app := &todo.App{DbPath: dbPath}
	err = app.NewApp(port)
	if err != nil {
		log.Fatal(err)
	}
	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"embed"
)

//go:embed web/dist
var svelteFiles embed.FS

func main() {
	checkArgs()
	initializeDB()
	initializeAuthDB()
	startHTTP()
}

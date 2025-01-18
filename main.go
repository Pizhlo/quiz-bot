package main

import (
	"flag"
	"quiz-mod/cmd/bot"
)

func main() {
	envFIle := flag.String("filename", ".env", "name of env file")
	path := flag.String("path", ".", "path to env file")
	configFile := flag.String("config file", "config.json", "name of config file")

	flag.Parse()

	bot.Start(*envFIle, *configFile, *path)
}

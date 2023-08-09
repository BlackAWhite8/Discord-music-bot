package main

import (
	"musicbot/bot"
	"os"
)

func main() {
	bot.RunBot(os.Args[1])
}

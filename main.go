package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"telegram-bot/bot"
)

func main() {
	telegramBot, err := bot.CreateBot()
	if err != nil {
		fmt.Printf("error creating News Bot: %v", err)
		os.Exit(1)
	}

	fmt.Println("NewsBot initialized successfully")

	signalsChannel := make(chan os.Signal, 1)
	signal.Notify(signalsChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	err = telegramBot.Run()
	if err != nil {
		fmt.Printf("error running bot: %v", err)
		os.Exit(1)
	}

	if err = telegramBot.StartGoRoutines(); err != nil {
		fmt.Printf("error starting goroutines: %v\n", err)
		os.Exit(1)
	}

	if err = telegramBot.StartHandlers(); err != nil {
		fmt.Printf("error starting handlers: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Waiting for sigterm")

	<-signalsChannel
	fmt.Println("exiting bot")
}

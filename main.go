package main

import (
	newsBot "bot-telegram/bot"
	"bot-telegram/scheduler"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	telegramNewsBot, err := newsBot.CreateNewsBot()
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

	err = telegramNewsBot.Run()
	if err != nil {
		fmt.Printf("error running bot: %v", err)
		os.Exit(1)
	}

	if err = telegramNewsBot.StartGoRoutines(); err != nil {
		fmt.Printf("error starting goroutines: %v\n", err)
		os.Exit(1)
	}

	if err = telegramNewsBot.StartHandlers(); err != nil {
		fmt.Printf("error starting handlers: %v\n", err)
		os.Exit(1)
	}

	newsScheduler := scheduler.NewNewsScheduler(telegramNewsBot, 3*time.Minute)
	newsScheduler.Run()

	fmt.Println("Waiting for sigterm")

	<-signalsChannel
	fmt.Println("exiting bot")
}

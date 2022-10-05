package main

import (
	"fmt"
	"os"
	"time"
)

const duration = time.Hour * 6

func ticker(action func() error) {
	ticker := time.NewTicker(duration)

	for ; true; <-ticker.C {
		run(action)
	}

}

func run(action func() error) {
	fmt.Println("Начало: " + time.Now().Format(time.RFC3339))

	err := action()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return

	}
	fmt.Println("Конец: " + time.Now().Format(time.RFC3339) + "\n")
}

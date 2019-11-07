package main

import "fmt"

func main() {
	messages := make(chan string, 2)

	messages <- "buffering"
	messages <- "channel"

	fmt.Println("messages1:", <-messages)
	fmt.Println("messages2:", <-messages)
}

package main

import (
	"fmt"

	"github.com/AaronRoethe/go-journal-client/message"
)

func main() {
	err := message.SendMessage()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Message sent successfully!")
}

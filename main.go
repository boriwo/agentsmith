package main

import (
	"log"
)

func main() {
	var err error
	sp, err := NewJSONSecretProvider("secrets.json")
	if err != nil {
		log.Fatalf("Error creating secret provider: %v", err)
	}
	ap := NewSimpleAnswerProvider()
	sa := NewSlackAgent(sp, ap)
	sa.LaunchSlack()
}

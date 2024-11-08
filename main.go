package main

import (
	"log"
)

var (
	agent Agent
)

func main() {
	var err error
	sessionMgr := NewSimpleSessionManager()
	kb := NewFileKnowledgeBase("kb")
	err = kb.Load()
	if err != nil {
		log.Fatalf("faild to load knowledge base: %v", err)
	}
	secretProvider, err := NewJSONSecretProvider("secrets.json")
	if err != nil {
		log.Fatalf("faild to create secret provider: %v", err)
	}
	answerProvider := NewStatefulAnswerProvider(kb)
	agent = NewSlackAgent(secretProvider, answerProvider, sessionMgr)
	agent.LaunchAgent()
}

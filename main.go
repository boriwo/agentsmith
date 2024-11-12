package main

import (
	"log"
)

const (
	DEFAULT_KNOWLEDGE_BASE_NAME = "kb"
	DEFAULT_EMBEDDING_BASE_NAME = "embeddings"
)

func main() {
	var err error
	sessionMgr := NewSimpleSessionManager()
	kb := NewFileKnowledgeBase(DEFAULT_KNOWLEDGE_BASE_NAME)
	err = kb.Load()
	if err != nil {
		log.Fatalf("faild to load knowledge base: %v", err)
	}
	secretProvider, err := NewJSONSecretProvider("secrets.json")
	if err != nil {
		log.Fatalf("faild to create secret provider: %v", err)
	}
	openaiHandler := NewOpenAIHandler(secretProvider)
	eb := NewEmbeddingBase(secretProvider, *openaiHandler, DEFAULT_EMBEDDING_BASE_NAME)
	err = eb.Load()
	if err != nil {
		log.Fatalf("faild to load embeddings base: %v", err)
	}
	err = eb.SyncEmbeddings(kb)
	if err != nil {
		log.Fatalf("faild to synchronize embeddings base: %v", err)
	}
	answerProvider := NewUberAnswerProvider(kb, eb, *openaiHandler)
	slackAgent := NewSlackAgent(secretProvider, answerProvider, sessionMgr)
	go slackAgent.LaunchAgent()
	webAgent := NewWebAgent(secretProvider, answerProvider, sessionMgr)
	webAgent.LaunchAgent()
}

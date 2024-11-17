/**
 * Copyright 2024 Boris Wolf
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	go webAgent.LaunchAgent()
	cliAgent := NewCliAgent(secretProvider, answerProvider, sessionMgr)
	cliAgent.LaunchAgent()
}

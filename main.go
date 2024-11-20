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
	"sync"
)

func main() {
	var err error
	sessionMgr := NewSimpleSessionManager()
	secretProvider, err := NewJSONSecretProvider("secrets.json")
	if err != nil {
		log.Fatalf("failed to create secret provider: %v", err)
	}
	configProvider, err := NewJSONConfigProvider("configs.json")
	if err != nil {
		log.Fatalf("failed to create config provider: %v", err)
	}
	openaiHandler := NewOpenAIHandler(secretProvider)
	if err != nil {
		log.Fatalf("failed to synchronize embeddings base: %v", err)
	}
	kbMgr, err := NewKnowledgeBaseManager(secretProvider, *openaiHandler)
	if err != nil {
		log.Fatalf("failed to load knowledge base: %v", err)
	}
	answerProvider := NewUberAnswerProvider(kbMgr, *openaiHandler)
	var wg sync.WaitGroup
	if configProvider.GetConfig("slackagent") == "yes" {
		slackAgent := NewSlackAgent(secretProvider, answerProvider, sessionMgr)
		wg.Add(1)
		go slackAgent.LaunchAgent(wg)
	}
	if configProvider.GetConfig("webagent") == "yes" {
		webAgent := NewWebAgent(configProvider, secretProvider, answerProvider, sessionMgr)
		wg.Add(1)
		go webAgent.LaunchAgent(wg)
	}
	if configProvider.GetConfig("cliagent") == "yes" {
		cliAgent := NewCliAgent(secretProvider, answerProvider, sessionMgr)
		wg.Add(1)
		go cliAgent.LaunchAgent(wg)
	}
	wg.Wait()
}

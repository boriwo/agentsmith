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
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var err error
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	sessionMgr := NewSimpleSessionManager()
	secretProvider, err := NewJSONSecretProvider("secrets.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to create secret provider")
	}
	configProvider, err := NewJSONConfigProvider("configs.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to create config provider")
	}
	if configProvider.GetConfig("loglevel") == "error" {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if configProvider.GetConfig("loglevel") == "warn" {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if configProvider.GetConfig("loglevel") == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	openaiHandler := NewOpenAIHandler(secretProvider)
	if err != nil {
		log.Error().Err(err).Msg("failed to synchronize embeddings base")
	}
	kbMgr, err := NewKnowledgeBaseManager(secretProvider, *openaiHandler)
	if err != nil {
		log.Error().Err(err).Msg("failed to load knowledge base")
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

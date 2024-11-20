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
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"sync"
)

type CliAgent struct {
	secretProvider SecretProvider
	answerProvider AnswerProvider
	sessionMgr     SessionManager
}

func NewCliAgent(secretProvider SecretProvider, answerProvider AnswerProvider, sessionManager SessionManager) Agent {
	cli := CliAgent{
		secretProvider: secretProvider,
		answerProvider: answerProvider,
		sessionMgr:     sessionManager,
	}
	return &cli
}

func (wa *CliAgent) LaunchAgent(wg sync.WaitGroup) {
	sessionId := ""
	log.Println("launching cli agent")
	fmt.Println("enter a question!")
	for {
		scanner := bufio.NewScanner(os.Stdin)
		question := ""
		if scanner.Scan() {
			question = scanner.Text()
			fmt.Printf("you entered: %s\n", question)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("error reading input: ", err)
		}
		if sessionId == "" {
			sessionId = wa.generateRandomString(12)
		}
		user := NewUser(sessionId, "CliUser", "CliUser")
		session := wa.sessionMgr.GetSession(user)
		answers, err := wa.answerProvider.GetAnswers(session, &Question{question})
		if err != nil {
			log.Println(err)
			return
		}
		for _, a := range answers {
			if a.Text != "" {
				fmt.Printf("%s\n", a.Text)
			}
			if a.Link != "" {
				fmt.Printf("%s\n", a.Link)
			}
			if a.ImageLink != "" {
				fmt.Printf("%s\n", a.ImageLink)
			}
		}
	}
	log.Println("stopping cli agent")
	wg.Done()
}

func (wa *CliAgent) generateRandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:n]
}

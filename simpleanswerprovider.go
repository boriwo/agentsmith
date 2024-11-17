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
	"fmt"
	"strings"
)

type SimpleAnswerProvider struct {
}

func NewSimpleAnswerProvider() AnswerProvider {
	return new(SimpleAnswerProvider)
}

func (sap *SimpleAnswerProvider) GetAnswers(session *UserSession, question *Question) ([]*Answer, error) {
	answer := new(Answer)
	if strings.Contains(question.Text, "hello") || strings.Contains(question.Text, "hi") {
		answer.Text = fmt.Sprintf("Hello %s", session.User.RealName)
	} else if strings.Contains(question.Text, "bye") {
		answer.Text = fmt.Sprintf("Good bye %s", session.User.RealName)
	} else {
		answer.Text = fmt.Sprintf("Sorry, I don't have the answers yet, %s", session.User.RealName)
	}
	answers := []*Answer{answer}
	session.LastQuestion = question
	session.LastAnswer = answers
	return answers, nil
}

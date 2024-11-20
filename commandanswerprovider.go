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
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	R_LIST_FACTS                 = "rlistfacts"
	R_LIST_KNOWLEDGE_BASES       = "rlistknowledgebases"
	R_NUM_FACTS                  = "rnumfacts"
	R_GET_FACT                   = "rgetfact"
	R_GET_CURRENT_KNOWLEDGE_BASE = "rgetcurrentknowledgebase"
	R_SET_CURRENT_KNOWLEDGE_BASE = "rsetcurrentknowledgebase"
	R_ADD_FACT                   = "raddfact"
	R_DELETE_FACT                = "rdeletefact"
)

type CommandAnswerProvider struct {
	kbm *KnowledeBaseManager
}

func NewCommandAnswerProvider(kbm *KnowledeBaseManager) AnswerProvider {
	answerProvider := CommandAnswerProvider{
		kbm,
	}
	return &answerProvider
}

func (sap *CommandAnswerProvider) GetAnswers(session *UserSession, question *Question) ([]*Answer, error) {
	answers := make([]*Answer, 0)
	answer := new(Answer)
	tokens := strings.Fields(question.Text)
	if len(tokens) > 0 && strings.HasPrefix(tokens[0], "<@") {
		tokens = tokens[1:]
	}
	if len(tokens) > 0 && tokens[0] == R_LIST_FACTS {
		for _, f := range sap.kbm.GetCurrentKnowledgeBase().ListFacts() {
			answer.Text += f.Name + "\n"
		}
		answers = append(answers, answer)
	} else if len(tokens) > 0 && tokens[0] == R_LIST_KNOWLEDGE_BASES {
		for _, name := range sap.kbm.ListBaseNames() {
			answer.Text += name + "\n"
		}
		answers = append(answers, answer)
	} else if len(tokens) > 0 && tokens[0] == R_GET_CURRENT_KNOWLEDGE_BASE {
		answer.Text += sap.kbm.GetCurrentBaseName() + "\n"
		answers = append(answers, answer)
	} else if len(tokens) > 0 && tokens[0] == R_SET_CURRENT_KNOWLEDGE_BASE {
		if len(tokens) < 2 {
			return nil, errors.New("missing parameter knowledge base name")
		}
		err := sap.kbm.SetCurrentBaseName(tokens[1])
		if err != nil {
			return nil, err
		}
		answer.Text += "set current knowledge base to " + tokens[1] + "\n"
		answers = append(answers, answer)
	} else if len(tokens) > 0 && tokens[0] == R_NUM_FACTS {
		answer.Text += fmt.Sprintf("%d", sap.kbm.GetCurrentKnowledgeBase().GetNumFacts())
		answers = append(answers, answer)
	} else if len(tokens) > 0 && tokens[0] == R_GET_FACT {
		if len(tokens) < 2 {
			return nil, errors.New("missing parameter fact name")
		}
		fact := sap.kbm.GetCurrentKnowledgeBase().GetFact(tokens[1])
		if fact == nil {
			return nil, errors.New("no fact by that name")
		} else {
			buf, _ := json.MarshalIndent(fact, "", "\t")
			answer.Text = string(buf)
			answers = append(answers, answer)
		}
	} else if len(tokens) > 0 && tokens[0] == R_ADD_FACT {
		if len(tokens) < 2 {
			return nil, errors.New("missing parameter fact name")
		}
		factName := tokens[1]
		if sap.kbm.GetCurrentKnowledgeBase().HasFact(factName) {
			return nil, errors.New("already have fact with name " + factName)
		}
		answer.Text += "adding new fact " + factName + ", please state a question for this fact!\n"
		answers = append(answers, answer)
		session.NewFact = new(Fact)
		session.NewFact.Name = strings.ToUpper(factName)
		session.NewFact.CreatedBy = session.User.Name
		session.NewFact.CreatedAt = fmt.Sprint(time.Now().Format(time.RFC3339))
		session.State = STATE_ADD_QUESTION
	} else if len(tokens) > 0 && tokens[0] == R_DELETE_FACT {
		if len(tokens) < 2 {
			return nil, errors.New("missing parameter fact name")
		}
		factName := tokens[1]
		if !sap.kbm.GetCurrentKnowledgeBase().HasFact(factName) {
			return nil, errors.New("no fact with name " + factName)
		}
		err := sap.kbm.GetCurrentKnowledgeBase().DeleteFact(factName)
		if err != nil {
			return nil, err
		}
		err = sap.kbm.GetCurrentEmbeddingsBase().SyncEmbeddings(sap.kbm.GetCurrentKnowledgeBase())
		if err != nil {
			return nil, err
		}
		err = sap.kbm.GetCurrentKnowledgeBase().Save()
		if err != nil {
			return nil, err
		}
	}
	session.LastQuestion = question
	session.LastAnswer = answers
	return answers, nil
}

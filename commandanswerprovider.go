package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	R_LIST_FACTS = "rlistfacts"
	R_NUM_FACTS  = "rnumfacts"
	R_GET_FACT   = "rgetfact"
)

type CommandAnswerProvider struct {
	kb KnowledeBaseProvider
	eb EmbeddingsBaseProvider
}

func NewCommandAnswerProvider(kb KnowledeBaseProvider, eb EmbeddingsBaseProvider) AnswerProvider {
	answerProvider := CommandAnswerProvider{
		kb,
		eb,
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
	if len(tokens) == 1 && tokens[0] == R_LIST_FACTS {
		for _, f := range sap.kb.ListFacts() {
			answer.Text += f.Name + "\n"
		}
		answers = append(answers, answer)
	} else if len(tokens) == 1 && tokens[0] == R_NUM_FACTS {
		answer.Text += fmt.Sprintf("%d", sap.kb.GetNumFacts())
		answers = append(answers, answer)
	} else if len(tokens) == 2 && tokens[0] == R_GET_FACT {
		fact := sap.kb.GetFact(tokens[1])
		if fact == nil {
			return nil, errors.New("no fact by that name")
		} else {
			buf, _ := json.MarshalIndent(fact, "", "\t")
			answer.Text = string(buf)
			answers = append(answers, answer)
		}
	}
	session.LastQuestion = question
	session.LastAnswer = answers
	return answers, nil
}
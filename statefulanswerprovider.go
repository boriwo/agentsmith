package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type StatefulAnswerProvider struct {
	kb KnowledeBaseProvider
}

func NewStatefulAnswerProvider(kb KnowledeBaseProvider) AnswerProvider {
	answerProvider := StatefulAnswerProvider{
		kb,
	}
	return &answerProvider
}

func (sap *StatefulAnswerProvider) GetAnswers(session *UserSession, question *Question) []*Answer {
	answers := make([]*Answer, 0)
	if session.State == STATE_QA {
		answer := new(Answer)
		tokens := strings.Fields(question.Text)
		if len(tokens) > 0 && strings.HasPrefix(tokens[0], "<@") {
			tokens = tokens[1:]
		}
		if len(tokens) == 1 && tokens[0] == "rlistfacts" {
			for _, f := range sap.kb.ListFacts() {
				answer.Text += f.Name + "\n"
			}
		} else if len(tokens) == 1 && tokens[0] == "rnumfacts" {
			answer.Text += fmt.Sprintf("%d", sap.kb.GetNumFacts())
		} else if len(tokens) == 2 && tokens[0] == "rgetfact" {
			fact := sap.kb.GetFact(tokens[1])
			if fact == nil {
				answer.Text = "no fact by that name"
			} else {
				buf, _ := json.MarshalIndent(fact, "", "\t")
				answer.Text = string(buf)
			}
		} else if strings.Contains(question.Text, "hello") || strings.Contains(question.Text, "hi") {
			answer.Text = fmt.Sprintf("Hello %s", session.User.RealName)
		} else if strings.Contains(question.Text, "bye") {
			answer.Text = fmt.Sprintf("Good bye %s", session.User.RealName)
		} else {
			answer.Text = fmt.Sprintf("Sorry, I don't have the answers yet, %s", session.User.RealName)
		}
		session.LastQuestion = question
		session.LastAnswer = answers
		session.State = STATE_QA
		answers = append(answers, answer)
	}
	return answers
}

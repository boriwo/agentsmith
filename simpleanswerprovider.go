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

func (sap *SimpleAnswerProvider) GetAnswers(session *UserSession, question *Question) []*Answer {
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
	session.State = STATE_QA
	return answers
}

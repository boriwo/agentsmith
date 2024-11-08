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

func (sap *SimpleAnswerProvider) GetAnswers(user *User, question *Question) []*Answer {
	a := new(Answer)
	if strings.Contains(question.Text, "hello") || strings.Contains(question.Text, "hi") {
		a.Text = fmt.Sprintf("Hello %s", user.RealName)
	} else if strings.Contains(question.Text, "bye") {
		a.Text = fmt.Sprintf("Good bye %s", user.RealName)
	} else {
		a.Text = fmt.Sprintf("Sorry, I don't have the answers yet, %s", user.RealName)
	}
	return []*Answer{a}
}

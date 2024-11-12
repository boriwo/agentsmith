package main

type ImageAnswerProvider struct {
	oai OpenAIHandler
}

func NewImageAnswerProvider(oai OpenAIHandler) AnswerProvider {
	answerProvider := ImageAnswerProvider{
		oai,
	}
	return &answerProvider
}

func (sap *ImageAnswerProvider) GetAnswers(session *UserSession, question *Question) ([]*Answer, error) {
	answers, err := sap.oai.GptGetImage(question)
	if err != nil {
		return nil, err
	}
	session.LastQuestion = question
	session.LastAnswer = answers
	return answers, nil
}

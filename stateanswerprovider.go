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

type StateAnswerProvider struct {
	kbm *KnowledeBaseManager
}

func NewStateAnswerProvider(kbm *KnowledeBaseManager) AnswerProvider {
	answerProvider := StateAnswerProvider{
		kbm,
	}
	return &answerProvider
}

func (sap *StateAnswerProvider) GetAnswers(session *UserSession, question *Question) ([]*Answer, error) {
	answers := make([]*Answer, 0)
	answer := new(Answer)
	if session.State == STATE_ADD_QUESTION {
		session.NewFact.Question = question.Text
		answer.Text += "please provide an answer to this question!\n"
		answers = append(answers, answer)
		session.State = STATE_ADD_ANSWER
	} else if session.State == STATE_ADD_ANSWER {
		if question.Text == "done" {
			err := sap.kbm.GetCurrentKnowledgeBase().AddFact(session.NewFact)
			if err != nil {
				answer.Text += "failed to add new fact " + session.NewFact.Name + " to knowledge base: " + err.Error() + "\n"
				answers = append(answers, answer)
			} else {
				err = sap.kbm.GetCurrentEmbeddingsBase().SyncEmbeddings(sap.kbm.GetCurrentKnowledgeBase())
				if err != nil {
					answer.Text += "failed to sync embeddings for new fact " + session.NewFact.Name + " for knowledge base: " + err.Error() + "\n"
					answers = append(answers, answer)
				}
				err = sap.kbm.GetCurrentKnowledgeBase().Save()
				if err != nil {
					answer.Text += "failed to save knwoledge base for new fact " + session.NewFact.Name + " for knowledge base: " + err.Error() + "\n"
					answers = append(answers, answer)
				} else {
					answer.Text += "added new fact " + session.NewFact.Name + " to knowledge base!\n"
					answers = append(answers, answer)
				}
			}
			session.State = STATE_QA
		} else {
			if session.NewFact.Answers == nil {
				session.NewFact.Answers = []string{}
			}
			session.NewFact.Answers = append(session.NewFact.Answers, question.Text)
			answer.Text += "please provide another answer to this question or type 'done' to finish and add fact!\n"
			answers = append(answers, answer)
		}
	} else {
		answer.Text += "unknown state " + session.State + ", reverting to default question/answer state\n"
		answers = append(answers, answer)
		session.State = STATE_QA
	}
	session.LastQuestion = question
	session.LastAnswer = answers
	return answers, nil
}

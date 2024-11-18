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
	"errors"
)

const (
	COMMAND_PLUGIN      = "COMMAND_PLUGIN"
	IMAGE_PLUGIN        = "IMAGE_PLUGIN"
	PARAM_TYPE_CONSTANT = "constant"
	PARAM_TYPE_PROMPT   = "prompt"
)

type PluginManager struct {
	kbm     *KnowledeBaseManager
	oai     OpenAIHandler
	plugins map[string]AnswerProvider
}

func NewPluginManger(kbm *KnowledeBaseManager, oai OpenAIHandler) *PluginManager {
	mgr := &PluginManager{
		kbm,
		oai,
		make(map[string]AnswerProvider),
	}
	mgr.plugins[COMMAND_PLUGIN] = NewCommandAnswerProvider(kbm)
	mgr.plugins[IMAGE_PLUGIN] = NewImageAnswerProvider(oai)
	return mgr
}

func (pm *PluginManager) GetAnswers(session *UserSession, question *Question, fact *Fact) ([]*Answer, error) {
	if question == nil || question.Text == "" {
		return nil, errors.New("missing question")
	}
	if fact == nil || fact.Plugin == "" {
		return nil, errors.New("missing fact or blank plugin")
	}
	answerProvider, ok := pm.plugins[fact.Plugin]
	if !ok {
		return nil, errors.New("unknown plugin " + fact.Plugin)
	}
	// rewrite question
	q := question.Text
	if len(fact.Params) > 0 {
		q = ""
		for idx, param := range fact.Params {
			if param.Type == PARAM_TYPE_CONSTANT {
				q += param.Value
				if idx < len(fact.Params)-1 {
					q += " "
				}
			} else if param.Type == PARAM_TYPE_PROMPT {
				answers, err := pm.oai.GptGetCompletions(&Question{param.Value + " " + question.Text})
				if err != nil {
					return nil, err
				}
				if len(answers) > 0 {
					q += answers[0].Text
					if idx < len(fact.Params)-1 {
						q += " "
					}
				}
			} else {
				return nil, errors.New("unsupported param type " + param.Type)
			}
		}
	}
	return answerProvider.GetAnswers(session, &Question{q})
}

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

import "sync"

const (
	STATE_QA = "STATE_QA"
)

type (
	User struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		RealName string `json:"realName"`
	}
	Parameter struct {
		Name             string `json:"name"`
		Value            string `json:"value"`
		Type             string `json:"type"`
		ExtractionPrompt string `json:"prompt"`
	}
	Fact struct {
		Name      string      `json:"name"`
		Question  string      `json:"question"` // quality text for generating embeddings
		Labels    []string    `json:"labels"`   // keywords for question
		Answers   []string    `json:"answers"`  // list of text based answers
		Links     []string    `json:"links"`    // list of http links
		Plugin    string      `json:"plugin"`   // optional plugin action
		Params    []Parameter `json:"params"`   // optional list of plugin params
		IsSystem  bool        `json:"isSystem"` // if true referring to a built in system command
		CreatedBy string      `json:"createdBy"`
		CreatedAt string      `json:"createdAt"`
	}
	Question struct {
		Text string
	}
	Answer struct {
		Text      string
		Link      string
		ImageLink string
		Score     float64
		Rank      int
	}
	UserSession struct {
		User         *User
		State        string
		LastQuestion *Question
		LastAnswer   []*Answer
	}
)

func NewQuestion(question string) *Question {
	q := Question{
		Text: question,
	}
	return &q
}

func NewAnswer(answer string) *Answer {
	a := Answer{
		Text: answer,
	}
	return &a
}

func (a *Answer) WithLink(link string) *Answer {
	a.Link = link
	return a
}

func (a *Answer) WithImageLink(link string) *Answer {
	a.ImageLink = link
	return a
}

func NewUser(id, name, realname string) *User {
	u := User{
		Id:       id,
		Name:     name,
		RealName: realname,
	}
	return &u
}

type Agent interface {
	LaunchAgent(wg sync.WaitGroup)
}

type AnswerProvider interface {
	GetAnswers(session *UserSession, question *Question) ([]*Answer, error)
}

type KnowledeBaseProvider interface {
	Load() error
	Save() error
	GetName() string
	GetFact(name string) *Fact
	GetNumFacts() int
	HasFact(name string) bool
	AddFact(fact *Fact) error
	DeleteFact(name string) error
	ListFacts() []*Fact
}

type EmbeddingsBaseProvider interface {
	Save() error
	Load() error
	GetName() string
	SyncEmbeddings(kb KnowledeBaseProvider) error
	RankEmbeddings(q *Embedding) (*EmbeddingsRanking, error)
	GetEmbedding(name string) *Embedding
	GetNumEmbeddings() int
	HasEmbedding(name string) bool
	AddEmbedding(embedding *Embedding) error
	DeleteEmbedding(name string) error
	ListEmbeddings() []*Embedding
}

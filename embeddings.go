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
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
)

type EmbeddingsBaseProvider interface {
	Save() error
	Load() error
	GetName() string
	SyncEmbeddings(kb KnowledeBaseProvider) error
	RankEmbeddings(q *Embedding) (*EmbeddingsRanking, error)
}

type FileEmbeddingsBase struct {
	name                 string
	filePath             string
	embeddingsByFactName map[string]*Embedding
	secretProvider       SecretProvider
	openaiHandler        OpenAIHandler
}

type EmbeddingsRanking struct {
	Embeddings []*Embedding
	Query      *Embedding
}

type Embedding struct {
	FactName      string    `json:"factId"`
	Source        string    `json:"source"`
	Link          string    `json:"link"`
	ModelId       string    `json:"modelId"`
	Embedding     []float64 `json:"embedding"`
	NumDimensions int       `json:"numDimensions"`
	Relevance     float64   `json:"relevance"`
}

func NewEmbedding(FactName, Source, Link, ModelId string) *Embedding {
	return &Embedding{
		FactName: FactName,
		Source:   Source,
		Link:     Link,
		ModelId:  ModelId,
	}
}

func (e *Embedding) WithEmbedding(vector []float64) *Embedding {
	e.Embedding = vector
	e.NumDimensions = len(vector)
	return e
}

func NewFileEmbeddingBase(secretProvider SecretProvider, openaiHandler OpenAIHandler, name string) EmbeddingsBaseProvider {
	eb := &FileEmbeddingsBase{
		embeddingsByFactName: make(map[string]*Embedding),
		secretProvider:       secretProvider,
		openaiHandler:        openaiHandler,
		name:                 name,
	}
	eb.filePath = filepath.Join("kb", "embeddings", eb.name+".json")
	return eb
}

func (e1 *Embedding) DotProd(e2 *Embedding) (float64, error) {
	if e1 == nil || e2 == nil {
		return 0, errors.New("empty embedding")
	}
	if len(e1.Embedding) != len(e2.Embedding) {
		return 0, errors.New("different dimensions")
	}
	if len(e1.Embedding) == 0 {
		return 0, errors.New("empty vectors")
	}
	/*if e1.ModelId != e2.ModelId {
		return 0, errors.New("different models")
	}*/
	p := float64(0.0)
	for idx, val := range e1.Embedding {
		p += val * e2.Embedding[idx]
	}
	return p, nil
}

func (e *Embedding) VecLen() (float64, error) {
	p, err := e.DotProd(e)
	if err != nil {
		return 0.0, err
	}
	return math.Sqrt(p), nil
}

func (e *Embedding) UpdateEmbedding(openaiHandler OpenAIHandler) error {
	if e.Source == "" {
		return errors.New("no source to embed")
	}
	log.Printf("updating embedding: %s\n", e.Source)
	newEmbedding, err := openaiHandler.GptGetEmbedding(&Question{e.Source})
	if err != nil {
		return err
	}
	e.Embedding = newEmbedding.Embedding
	e.NumDimensions = newEmbedding.NumDimensions
	e.ModelId = ""
	return nil
}

func (e *Embedding) Clone() *Embedding {
	return &Embedding{
		FactName:      e.FactName,
		Source:        e.Source,
		Link:          e.Link,
		ModelId:       e.ModelId,
		Embedding:     e.Embedding,
		NumDimensions: len(e.Embedding),
		Relevance:     e.Relevance,
	}
}

func (eb *FileEmbeddingsBase) GetName() string {
	return eb.name
}

func (eb *FileEmbeddingsBase) Save() error {
	a := make([]*Embedding, len(eb.embeddingsByFactName))
	idx := 0
	for _, e := range eb.embeddingsByFactName {
		a[idx] = e
		idx++
	}
	sort.SliceStable(a, func(i, j int) bool {
		if a[i].FactName < a[j].FactName {
			return true
		} else {
			return false
		}
	})
	buf, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(eb.filePath, buf, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (eb *FileEmbeddingsBase) Load() error {
	data, err := os.ReadFile(eb.filePath)
	if err != nil {
		return err
	}
	a := make([]*Embedding, 0)
	err = json.Unmarshal(data, &a)
	if err != nil {
		return err
	}
	eb.embeddingsByFactName = make(map[string]*Embedding, 0)
	for _, e := range a {
		if e.FactName == "" {
			return errors.New("invalid fact")
		}
		eb.embeddingsByFactName[e.FactName] = e
	}
	return nil
}

func (eb *FileEmbeddingsBase) SyncEmbeddings(kb KnowledeBaseProvider) error {
	var err error
	if kb == nil {
		return errors.New("no knowedge base")
	}
	changed := false
	for _, fact := range kb.ListFacts() {
		emb, ok := eb.embeddingsByFactName[fact.Name]
		if ok {
			if emb.Source != fact.Question && fact.Question != "" {
				emb.Source = fact.Question
				err = emb.UpdateEmbedding(eb.openaiHandler)
				changed = true
				if err != nil {
					return err
				}
			}
			if len(emb.Embedding) == 0 || emb.NumDimensions == 0 || len(emb.Embedding) != emb.NumDimensions {
				emb.Source = fact.Question
				err = emb.UpdateEmbedding(eb.openaiHandler)
				changed = true
				if err != nil {
					return err
				}
			}
		} else {
			if fact.Question != "" {
				eb.embeddingsByFactName[fact.Name] = NewEmbedding(fact.Name, fact.Question, "", GPT_CURRENT_MODEL)
				err = eb.embeddingsByFactName[fact.Name].UpdateEmbedding(eb.openaiHandler)
				changed = true
				if err != nil {
					return err
				}
			}
		}
	}
	for factName, _ := range eb.embeddingsByFactName {
		if !kb.HasFact(factName) {
			delete(eb.embeddingsByFactName, factName)
			changed = true
		}
	}
	if changed {
		err = eb.Save()
		if err != nil {
			return err
		}
	}
	return nil
}

func (eb *FileEmbeddingsBase) RankEmbeddings(q *Embedding) (*EmbeddingsRanking, error) {
	er := &EmbeddingsRanking{
		Embeddings: make([]*Embedding, len(eb.embeddingsByFactName)),
		Query:      q,
	}
	if len(eb.embeddingsByFactName) == 0 {
		return er, errors.New("no embeddings")
	}
	if q.Source == "" {
		return er, errors.New("no query")
	}
	if len(q.Embedding) == 0 {
		return er, errors.New("no query embedding")
	}
	idx := 0
	var err error
	for _, e := range eb.embeddingsByFactName {
		er.Embeddings[idx] = e.Clone()
		er.Embeddings[idx].Relevance, err = er.Embeddings[idx].DotProd(q)
		if err != nil {
			return er, err
		}
		idx++
	}
	sort.SliceStable(er.Embeddings, func(i, j int) bool {
		dpi, ierr := er.Embeddings[i].DotProd(q)
		if ierr != nil {
			err = ierr
		}
		dpj, ierr := er.Embeddings[j].DotProd(q)
		if ierr != nil {
			err = ierr
		}
		if dpi > dpj {
			return true
		} else {
			return false
		}
	})
	return er, nil
}

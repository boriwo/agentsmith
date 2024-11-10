package main

import (
	"encoding/json"
	"errors"
	"math"
	"os"
	"sort"
)

const ()

type EmbeddingsBase struct {
	Name               string
	EmbeddingsByFactId map[string]*Embedding
}

type EmbeddingsRanking struct {
	Embeddings []*Embedding
	Query      *Embedding
}

type Embedding struct {
	FactId        string    `json:"factId"`
	Source        string    `json:"source"`
	Link          string    `json:"link"`
	ModelId       string    `json:"modelId"`
	Embedding     []float64 `json:"embedding"`
	NumDimensions int       `json:"numDimensions"`
	Relevance     float64   `json:"relevance"`
}

func NewEmbedding(FactId, Source, Link, ModelId string) *Embedding {
	return &Embedding{
		FactId:  FactId,
		Source:  Source,
		Link:    Link,
		ModelId: ModelId,
	}
}

func (e *Embedding) WithEmbedding(vector []float64) *Embedding {
	e.Embedding = vector
	e.NumDimensions = len(vector)
	return e
}

func NewEmbeddingBase(name string) *EmbeddingsBase {
	return &EmbeddingsBase{
		EmbeddingsByFactId: make(map[string]*Embedding),
		Name:               name,
	}
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
	if e1.ModelId != e2.ModelId {
		return 0, errors.New("different models")
	}
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

func (e *Embedding) Clone() *Embedding {
	return &Embedding{
		FactId:        e.FactId,
		Source:        e.Source,
		Link:          e.Link,
		ModelId:       e.ModelId,
		Embedding:     e.Embedding,
		NumDimensions: len(e.Embedding),
		Relevance:     e.Relevance,
	}
}

func (eb *EmbeddingsBase) StoreToDisk() error {
	a := make([]*Embedding, len(eb.EmbeddingsByFactId))
	idx := 0
	for _, e := range eb.EmbeddingsByFactId {
		a[idx] = e
		idx++
	}
	sort.SliceStable(a, func(i, j int) bool {
		if a[i].FactId < a[j].FactId {
			return true
		} else {
			return false
		}
	})
	buf, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(eb.Name+".json", buf, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (eb *EmbeddingsBase) LoadFromDisk() error {
	data, err := os.ReadFile(eb.Name + ".json")
	if err != nil {
		return err
	}
	a := make([]*Embedding, 0)
	err = json.Unmarshal(data, &a)
	if err != nil {
		return err
	}
	eb.EmbeddingsByFactId = make(map[string]*Embedding, 0)
	for _, e := range a {
		if e.FactId == "" {
			return errors.New("invalid fact")
		}
		eb.EmbeddingsByFactId[e.FactId] = e
	}
	return nil
}

func (e *Embedding) UpdateVector() error {
	if e.Source == "" {
		return errors.New("no source to embed")
	}
	//TODO: update embedding here
	return nil
}

func (eb *EmbeddingsBase) SyncEmbeddings(kb KnowledeBaseProvider) error {
	var err error
	changed := false
	for _, fact := range kb.ListFacts() {
		emb, ok := eb.EmbeddingsByFactId[fact.Name]
		if ok {
			if emb.Source != fact.Question && fact.Question != "" {
				emb.Source = fact.Question
				err = emb.UpdateVector()
				changed = true
				if err != nil {
					return err
				}
			}
			if len(emb.Embedding) == 0 || emb.NumDimensions == 0 || len(emb.Embedding) != emb.NumDimensions {
				emb.Source = fact.Question
				err = emb.UpdateVector()
				changed = true
				if err != nil {
					return err
				}
			}
		} else {
			if fact.Question != "" {
				eb.EmbeddingsByFactId[fact.Name] = NewEmbedding(fact.Name, fact.Question, "", GPT_CURRENT_MODEL)
				err = eb.EmbeddingsByFactId[fact.Name].UpdateVector()
				changed = true
				if err != nil {
					return err
				}
			}
		}
	}
	for factId, _ := range eb.EmbeddingsByFactId {
		if !kb.HasFact(factId) {
			delete(eb.EmbeddingsByFactId, factId)
			changed = true
		}
	}
	if changed {
		err = eb.StoreToDisk()
		if err != nil {
			return err
		}
	}
	return nil
}

func (eb *EmbeddingsBase) RankEmbeddings(q *Embedding) (*EmbeddingsRanking, error) {
	er := &EmbeddingsRanking{
		Embeddings: make([]*Embedding, len(eb.EmbeddingsByFactId)),
		Query:      q,
	}
	if len(eb.EmbeddingsByFactId) == 0 {
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
	for _, e := range eb.EmbeddingsByFactId {
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

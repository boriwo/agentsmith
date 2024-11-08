package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type FileKnowledgeBase struct {
	sync.Mutex
	Name  string           `json:"name"`
	facts map[string]*Fact `json:"facts"`
}

func NewFileKnowledgeBase(name string) KnowledeBaseProvider {
	kb := FileKnowledgeBase{
		Name:  name,
		facts: make(map[string]*Fact, 0),
	}
	return &kb
}

func (fkb *FileKnowledgeBase) Load(name string) error {
	fkb.Lock()
	defer fkb.Unlock()
	jsonData, err := os.ReadFile(fkb.Name + ".json")
	if err != nil {
		return err
	}
	var facts []*Fact
	err = json.Unmarshal(jsonData, &facts)
	if err != nil {
		return err
	}
	fkb.facts = make(map[string]*Fact)
	for _, f := range facts {
		fkb.facts[f.Name] = f
	}
	return nil
}

func (fkb *FileKnowledgeBase) Save(name string) error {
	fkb.Lock()
	defer fkb.Unlock()
	var facts []*Fact
	for _, f := range fkb.facts {
		facts = append(facts, f)
	}
	jsonData, err := json.MarshalIndent(facts, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(fkb.Name+".json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (fkb *FileKnowledgeBase) GetName() string {
	return fkb.Name
}

func (fkb *FileKnowledgeBase) GetFact(name string) *Fact {
	fkb.Lock()
	defer fkb.Unlock()
	return fkb.facts[name]
}

func (fkb *FileKnowledgeBase) GetNumFacts() int {
	fkb.Lock()
	defer fkb.Unlock()
	return len(fkb.facts)
}

func (fkb *FileKnowledgeBase) HasFact(name string) bool {
	fkb.Lock()
	defer fkb.Unlock()
	_, ok := fkb.facts[name]
	return ok
}

func (fkb *FileKnowledgeBase) AddFact(fact *Fact) error {
	fkb.Lock()
	defer fkb.Unlock()
	if fact == nil {
		return errors.New("empty fact")
	}
	if fact.Name == "" {
		return errors.New("fact needs name")
	}
	_, ok := fkb.facts[fact.Name]
	if ok {
		return errors.New("fact already exists")
	}
	fkb.facts[fact.Name] = fact
	return nil
}

func (fkb *FileKnowledgeBase) DeleteFAct(name string) error {
	fkb.Lock()
	defer fkb.Unlock()
	_, ok := fkb.facts[name]
	if !ok {
		return fmt.Errorf("no fact with name %s exists", name)
	}
	delete(fkb.facts, name)
	return nil
}

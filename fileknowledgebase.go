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
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileKnowledgeBase struct {
	sync.Mutex
	name     string
	filePath string
	facts    map[string]*Fact
}

func NewFileKnowledgeBase(name string) KnowledeBaseProvider {
	kb := FileKnowledgeBase{
		name:  name,
		facts: make(map[string]*Fact, 0),
	}
	kb.filePath = filepath.Join("kb", "facts", kb.name+".json")
	return &kb
}

func (fkb *FileKnowledgeBase) Load() error {
	fkb.Lock()
	defer fkb.Unlock()
	jsonData, err := os.ReadFile(fkb.filePath)
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

func (fkb *FileKnowledgeBase) Save() error {
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
	err = os.WriteFile(fkb.filePath, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (fkb *FileKnowledgeBase) GetName() string {
	return fkb.name
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

func (fkb *FileKnowledgeBase) DeleteFact(name string) error {
	fkb.Lock()
	defer fkb.Unlock()
	_, ok := fkb.facts[name]
	if !ok {
		return fmt.Errorf("no fact with name %s exists", name)
	}
	delete(fkb.facts, name)
	return nil
}

func (fkb *FileKnowledgeBase) ListFacts() []*Fact {
	fkb.Lock()
	defer fkb.Unlock()
	allFacts := make([]*Fact, 0)
	for _, f := range fkb.facts {
		allFacts = append(allFacts, f)
	}
	return allFacts
}

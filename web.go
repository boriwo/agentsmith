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
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type WebAgent struct {
	secretProvider SecretProvider
	configProvider ConfigProvider
	answerProvider AnswerProvider
	sessionMgr     SessionManager
}

func NewWebAgent(configProvider ConfigProvider, secretProvider SecretProvider, answerProvider AnswerProvider, sessionManager SessionManager) Agent {
	wa := WebAgent{
		configProvider: configProvider,
		secretProvider: secretProvider,
		answerProvider: answerProvider,
		sessionMgr:     sessionManager,
	}
	return &wa
}

func (wa *WebAgent) LaunchAgent(wg sync.WaitGroup) {
	r := mux.NewRouter().StrictSlash(true)
	fs := http.FileServer(http.Dir("./web"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/agentsmith", wa.getHandler).Methods("GET")
	r.HandleFunc("/agentsmith", wa.postHandler).Methods("POST")
	log.Println("launching web agent")
	http.ListenAndServe(wa.configProvider.GetConfig("webport"), r)
	log.Println("stopping web agent")
	wg.Done()
}

func (wa *WebAgent) generateRandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:n]
}

func (wa *WebAgent) getHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := wa.generateRandomString(12)
	data := map[string]string{
		"Question":     "",
		"LastQuestion": "",
		"Answer":       "",
		"AnswerTitle":  "",
		"AnswerLink":   "",
		"AnswerImage":  "",
		"SessionId":    sessionId,
		"Error":        "",
	}
	tmpl, err := template.ParseFiles("web/form.html")
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
}

func (wa *WebAgent) postHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	question := r.FormValue("question")
	sessionId := r.FormValue("sessionId")
	if sessionId == "" {
		sessionId = wa.generateRandomString(12)
	}
	user := NewUser(sessionId, "WebUser", "WebUser")
	session := wa.sessionMgr.GetSession(user)
	data := map[string]string{
		"Question":     "",
		"Answer":       "",
		"AnswerLink":   "",
		"AnswerImage":  "",
		"AnswerTitle":  "Answer",
		"SessionId":    sessionId,
		"LastQuestion": "",
		"Error":        "",
	}
	answers, err := wa.answerProvider.GetAnswers(session, &Question{question})
	if err != nil {
		data["Error"] = err.Error()
	} else {
		for _, a := range answers {
			data["Answer"] += a.Text
			data["AnswerLink"] += a.Link
			data["AnswerImage"] += a.ImageLink
		}
	}
	if session.LastQuestion != nil {
		data["LastQuestion"] = session.LastQuestion.Text
	}
	tmpl, err := template.ParseFiles("web/form.html")
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
}

package main

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type WebAgent struct {
	secretProvider SecretProvider
	answerProvider AnswerProvider
	sessionMgr     SessionManager
}

func NewWebAgent(secretProvider SecretProvider, answerProvider AnswerProvider, sessionManager SessionManager) Agent {
	wa := WebAgent{
		secretProvider: secretProvider,
		answerProvider: answerProvider,
		sessionMgr:     sessionManager,
	}
	return &wa
}

func (wa *WebAgent) LaunchAgent() {
	r := mux.NewRouter().StrictSlash(true)
	fs := http.FileServer(http.Dir("./web"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/agentsmith", wa.getHandler).Methods("GET")
	r.HandleFunc("/agentsmith", wa.postHandler).Methods("POST")
	http.ListenAndServe(wa.secretProvider.GetSecret("webport"), r)
}

func (wa *WebAgent) generateRandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:n]
}

func (wa *WebAgent) getHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := wa.generateRandomString(12)
	data := map[string]string{
		"Question":    "",
		"Answer":      "",
		"AnswerTitle": "",
		"AnswerLink":  "",
		"AnswerImage": "",
		"SessionId":   sessionId,
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
	answers, err := wa.answerProvider.GetAnswers(session, &Question{question})
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	data := map[string]string{
		"Question":    question,
		"Answer":      "",
		"AnswerLink":  "",
		"AnswerImage": "",
		"AnswerTitle": "Answer",
		"SessionId":   sessionId,
	}
	for _, a := range answers {
		data["Answer"] += a.Text
		data["AnswerLink"] += a.Link
		data["AnswerImage"] += a.ImageLink
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

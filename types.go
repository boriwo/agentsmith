package main

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
	LaunchAgent()
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
	AddFact(*Fact) error
	DeleteFAct(name string) error
	ListFacts() []*Fact
}

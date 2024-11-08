package main

type (
	User struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		RealName string `json:"real_name"`
	}
	Fact struct {
		Name           string   `json:"name"`
		Question       string   `json:"question"`        // quality text for generating embeddings
		Labels         []string `json:"labels"`          // keywords for question
		Answers        []string `json:"answers"`         // list of text based answers
		Links          []string `json:"links"`           // list of http links
		Plugin         string   `json:"plugin"`          // optional plugin action
		ParamSignature []string `json:"param_signature"` // optional list of parameters to be found in question to be passed to plugin
		ParamPrompt    string   `json:"param_prompt"`    // optional prompt for extracting parameter key value pairs to be passed to plugin
		CreatedBy      string   `json:"created_by"`
		CreatedAt      string   `json:"created_at"`
	}
	Question struct {
		Text string
	}
	Answer struct {
		Text  string
		Score float64
		Rank  int
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

func NewUser(id, name, realname string) *User {
	u := User{
		Id:       id,
		Name:     name,
		RealName: realname,
	}
	return &u
}

type AnswerProvider interface {
	GetAnswers(user *User, question *Question) []*Answer
}

type KnowledeBaseProvider interface {
	Load(name string) error
	Save(name string) error
	GetName() string
	GetFact(name string) *Fact
	GetNumFacts() int
	HasFact(name string) bool
	AddFact(*Fact) error
	DeleteFAct(name string) error
}

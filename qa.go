// Question Answering System implemention.

package Posger

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"io"
	"bytes"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"fmt"
	"net/url"
)

const (
	_ANSWERING_SERVER = "/answer?question=%s&question_id=%s"
)

func registeQuestionApi(router *mux.Router) {
	router.HandleFunc("/answer", processQuestion).Methods("GET")
	router.HandleFunc("/answer/{questionId}", getAnswer).Methods("GET")
	router.HandleFunc("/answer/{questionId}", alterAnswer).Methods("POST")
}

// The client requests the server, providing a question parameter with question
// The answer response by
func processQuestion(w http.ResponseWriter, r *http.Request) {
	if q := r.URL.Query()["question"]; len(q) == 0 {
		// return error msg: question parameter is required
	} else {
		questionId := ""
		// TODO: if the question content has proceeded before, we can directly return the cache answer in database
		questions := SelectQuestion(map[string]interface{}{"question": q[0]})
		if len(questions) == 0 {
			questionId = uuid.Must(uuid.NewV4()).String()
			// request the bottle server with question and questionId
			form := url.Values{}
			form.Add("question", q[0])
			form.Add("question_id", questionId)
			fmt.Println(fmt.Sprintf(_ANSWERING_SERVER, q[0], questionId))
			_, err := http.Get("http://0.0.0.0:8081/answer?" + form.Encode())
			if err != nil {
				Logger.Println("Request question server error")
				d, _ := json.Marshal(struct {
					Msg string			`json:"error"`
					QuestionId	string  `json:"question_id"`
				}{"Answering server error...", questionId})
				io.Copy(w, bytes.NewReader(d))
				return
			}
			go AddQuestion(Question{questionId, q[0], "", nil, nil, nil})
			// return message, if server response successfully, waiting for answer
			// return error, if server
		} else {
			questionId = questions[0].QuestionId
		}
		d, _ := json.Marshal(struct {
			Msg string			`json:"msg"`
			QuestionId	string	`json:"question_id"`
		}{"Nice request!", questionId})
		io.Copy(w, bytes.NewReader(d))

		//question := q[0]
		// record question Id, return successful and allocated questionId, processing message.
		// request python bottle server to get question's answer
		// next js client to request server by question id to get answer
		// When
	}
}

// Get Answer by questionId predefined in system.
func getAnswer(w http.ResponseWriter, r *http.Request) {
	questionId := mux.Vars(r)["questionId"]
	if qs := SelectQuestion(map[string]interface{}{"questionid": questionId}); len(qs) == 0{
		// return questionId is wrong message
	} else {
		question := qs[0]
		if question.Answer == "" {
			// the answer still is proceeding
			d, _ := json.Marshal(struct {
				Msg string		`json:"error"`
				Question string	`json:"question"`
				Answer	string	`json:"answer"`
			}{"Answer Completed!", question.Question, question.Answer})
			io.Copy(w, bytes.NewReader(d))
		} else {
			// return the proceeded question / answer message
			d, _ := json.Marshal(struct {
				Msg string		`json:"msg"`
				Question	`json:"question"`
				Answer	string	`json:"answer"`
			}{"Answer Completed!", question, question.Answer})
			io.Copy(w, bytes.NewReader(d))
		}
	}
}

// Add answer and passages by questionId, request by ZeusKnows Server
// POST
func alterAnswer(w http.ResponseWriter, r *http.Request) {
	questionId := mux.Vars(r)["questionId"]
	body, _ := ioutil.ReadAll(r.Body)
	info := &struct{
		Answer		string		`json:"answer"`
		Scores		[]float32	`json:"scores"`
		Answers		[]string	`json:"answers"`
		Passages	[]string	`json:"passages"`
	}{}
	json.Unmarshal(body, &info)
	// answer and passage is saved in body by json format.
	if qs := SelectQuestion(map[string]interface{}{"questionid": questionId}); len(qs) == 0{
		// return questionId is wrong message
		io.Copy(w, bytes.NewReader([]byte(`{"error": "question_id didn't exist!"}`)))
	} else {
		question := qs[0]
		go SetQuestionAnswer(Question{question.QuestionId, question.Question, info.Answer, info.Scores, info.Answers, info.Passages})
		io.Copy(w, bytes.NewReader([]byte(`{"msg": "ok"}"`)))
		Logger.Println(question.Question, ": ", info.Answer)
	}
}






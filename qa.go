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
)

const (
	_ANSWERING_SERVER = "0.0.0.0:8081/answer?question="
)

func registeQuestionApi(router *mux.Router) {
	router.HandleFunc("/", processQuestion)
	router.HandleFunc("/answer/{questionId}", getAnswer)
}

// The client requests the server, providing a question parameter with question
// The answer response by
func processQuestion(w http.ResponseWriter, r *http.Request) {
	if q := r.URL.Query()["question"]; len(q) == 0 {
		// return error msg: question parameter is required
	} else {
		// TODO: if the question content has proceeded before, we can directly return the cache answer in database
		questions := SelectQuestion(map[string]interface{}{"Question": q})
		if len(questions) == 1 {
			d, _ := json.Marshal(struct {
				Msg string			`json:"msg"`
				QuestionId	string	`json:"question_id"`
			}{"Nice request!", questions[0].QuestionId})
			io.Copy(w, bytes.NewReader(d))
		} else {
			// request the bottle server with question and questionId
			go AddQuestion(Question{uuid.Must(uuid.NewV4()).String(), q[0], "", nil})
			// return message, if server response successfully, waiting for answer
			// return error, if server
			resp, err := http.Get(_ANSWERING_SERVER + q[0])
			if err != nil {
				Logger.Fatalln("Request question server error")
				d, _ := json.Marshal(struct {
					Msg string			`json:"msg"`
					QuestionId	string	`json:"question_id"`
				}{"Nice request!", questions[0].QuestionId})
				io.Copy(w, bytes.NewReader(d))
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

		}

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
		} else {
			// return the proceeded question / answer message
		}
	}
}






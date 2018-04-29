// Question Answering System implemention.

package Posger

import (
	"github.com/gorilla/mux"
	"net/http"
)

func registeQuestionApi(router *mux.Router) {
	router.HandleFunc("/", processQuestion)
	router.HandleFunc("/answer/{questionId}", getAnswer)
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

// The client requests the server, providing a question parameter with question
// The answer response by
func processQuestion(w http.ResponseWriter, r *http.Request) {
	if q := r.URL.Query()["question"]; len(q) == 0 {
		// return error msg: question parameter is required
	} else {
		// TODO: if the question content has proceeded before, we can directly return the cache answer in database
		//question := q[0]
		// record question Id, return successful and allocated questionId, processing message.
		// request python bottle server to get question's answer
		// next js client to request server by question id to get answer
		// When
	}
}




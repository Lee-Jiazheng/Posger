package Posger

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"io"
	"bytes"
	"encoding/json"
)

func registeAjaxApi(router *mux.Router) {
	router.HandleFunc("/login", login).Methods("GET")
	router.HandleFunc("/logout", logout).Methods("GET")
	registeSummaryApi(router.PathPrefix("/digest").Subrouter())
}

// API's required login decorator.
// The f's parameters is filter, return string is json string.
func RequireLoginApi(w http.ResponseWriter, r *http.Request, f func(string) []byte) {
	user, err := r.Cookie("user")
	if err != nil {
		io.Copy(w, bytes.NewReader([]byte(`{"error": "Please login first"}`)))
	} else {
		io.Copy(w, bytes.NewReader(f(user.Value)))
	}
}

// Check whether the username and password match.
func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	users := SelectUser(map[string]interface{}{"username": r.Form.Get("username"), "password": r.Form.Get("password")})
	if len(users) == 1 {
		// Later, we will marsh a better user info cookie.
		http.SetCookie(w, &http.Cookie{
			Name: "user",
			Value: r.Form.Get("username"),	// user struct json
			Path: "/",
			Expires: time.Now().AddDate(0, 1, 0),
			MaxAge: 86400,	// 100 hours' validate time
		})
		d, _ := json.Marshal(struct {
			Msg string
		}{"login successfully!"})
		io.Copy(w, bytes.NewReader(d))
	} else {
		d, _ := json.Marshal(struct {
			Error string	`json:"error"`
		}{"username or password incorrect!"})
		io.Copy(w, bytes.NewReader(d))
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	// clear the cookie of "user"
	http.SetCookie(w, &http.Cookie{Name: "user", Path: "/", MaxAge: -1})
	http.Redirect(w, r, "/index", http.StatusFound)
}
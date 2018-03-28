package Posger

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"html/template"
	"io"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"fmt"
)

var (
	hostAddress string
)

func init() {
	hostAddress = "127.0.0.1:8080"
}

func RunServer() {
	router := mux.NewRouter()
	router.HandleFunc("/index", indexHandler)
	router.Handle("/", http.RedirectHandler("/index", 301))

	//train_router := router.PathPrefix("/train").Subrouter()
	
	summary_router := router.PathPrefix("/summarize").Subrouter()
	summary_router.HandleFunc("/upload", uploadView).Methods("GET")
	summary_router.HandleFunc("/upload", uploadPaper).Methods("POST")
	summary_router.HandleFunc("/poster/{paperId}", summarizePaper).Methods("GET")
	summary_router.HandleFunc("/info/{paperId}", articleInfo).Methods("GET")

	registerOauth2App(router.PathPrefix("/oauth2").Subrouter())
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))


	srv := &http.Server{
		Handler:      router,
		Addr:         hostAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Server listened at " + "http://127.0.0.1:8080")
	log.Fatal(srv.ListenAndServe())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/views/index.html")
	t.Execute(w, nil)
}

func uploadView(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/views/upload.html")
	t.Execute(w, nil)
}

func uploadPaper(w http.ResponseWriter, r *http.Request) {
	up_file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer up_file.Close()
	new_f, _ := os.Create("static/articles/" + uuid.Must(uuid.NewV4()).String())
	defer new_f.Close()
	io.Copy(new_f, up_file)
	http.Redirect(w, r, "/summarize/paper/1", http.StatusFound)
}

func summarizePaper(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//vars["paperId"]
	t, _ := template.ParseFiles("static/views/sum_template.html")
	t.Execute(w, nil)
}

func articleInfo(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//vars["paperId"]
	article, err := NewJsonArticle("static/articles/大数据时代我国企业财务共享中心的优化.pdf")
	if err != nil {
		//404
	}
	if d, err := json.Marshal(article); err != nil {
		fmt.Print(err)
	} else {
		io.Copy(w, bytes.NewReader(d))
	}
}
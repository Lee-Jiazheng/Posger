package Posger

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"html/template"
	"io"
	"encoding/json"
	"net/http"
	"os"
	"time"
	"fmt"
	"log"
)

var (
	hostAddress string
)

func init() {
	hostAddress = "0.0.0.0:8080"
}

func RunServer() {

	router := mux.NewRouter()
	// Basic View Config
	router.HandleFunc("/index", indexView).Methods("GET")
	router.HandleFunc("/digest/{paperId}", digestView).Methods("GET")
	router.HandleFunc("/q-a", questionView).Methods("GET")
	router.HandleFunc("/userinfo/{userId}", infoView).Methods("GET")
	router.HandleFunc("/contact", contactView).Methods("GET")

	//train_router := router.PathPrefix("/train").Subrouter()
	
	summary_router := router.PathPrefix("/digests").Subrouter()
	summary_router.HandleFunc("/upload", uploadView).Methods("GET")
	summary_router.HandleFunc("/upload", uploadPaper).Methods("POST")
	summary_router.HandleFunc("/poster/{paperId}", summarizePaper).Methods("GET")
	summary_router.HandleFunc("/info/{paperId}", articleInfo).Methods("GET")

	registeOauth2App(router.PathPrefix("/oauth2").Subrouter())
	registeAjaxApi(router.PathPrefix("/api").Subrouter())

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Exception View Config
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/index", http.StatusFound)
	})
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Logger.Println("Resource Not Found: ", r.URL.Path)
		t, _ := template.ParseFiles("static/views/404.html", "static/views/ref.html")
		t.Execute(w, checkLoginUser(r, "404"))
	})

	srv := &http.Server{
		Handler:      router,
		Addr:         hostAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Server listened at " + "http://" + hostAddress)
	Logger.Println(srv.ListenAndServe())
}

// return index web page
func indexView(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/views/index.html", "static/views/ref.html")
	t.Execute(w, checkLoginUser(r, "index"))
}

// return digest web page, methods: get
func digestView(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/views/digest.html", "static/views/ref.html")
	t.Execute(w, checkLoginUser(r, "digest"))
}

// return qustion-answering web page, methods: get
func questionView(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/views/answer.html", "static/views/ref.html")
	t.Execute(w, checkLoginUser(r, "question"))
}

func convert2time(timestamp int32) string {
	fmt.Println(timestamp)
	return time.Unix(int64(timestamp), 0).Format("2006-01-02 15:04:05")
}

// If user has already logged in, it will show its personal page.
// Otherwise, it will redirect to the index page.
func infoView(w http.ResponseWriter, r *http.Request) {
	if username := isLogin(r); username == "anonymous" {
		http.Redirect(w, r, "/index", http.StatusFound)
	} else {
		// The template Name must be "user.html" as base template
		// Must define funcs before parse files.
		t, _ := template.New("user.html").Funcs(template.FuncMap{"showTime": func(ts int32) string {
			return time.Unix(int64(ts), 0).Format("2006-01-02 15:04:05")
		}}).ParseFiles("static/views/user.html", "static/views/ref.html")

		if users := SelectUser(map[string]interface{}{"username": username, "userid": mux.Vars(r)["userId"]}); len(users) != 0 {
			if err := t.Execute(w, struct {
				User	*User
				Papers	[]Paper
				PageType	string
			}{&users[0], SelectPaper(map[string]interface{}{"owner": mux.Vars(r)["userId"]}), "info"}); err != nil {
				fmt.Println(err)
			}
		} else {
			http.Redirect(w, r, "/index", http.StatusFound)
		}
	}
}

// Contact web page, introduce the developer.
func contactView(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/views/contactus.html", "static/views/ref.html")
	t.Execute(w, checkLoginUser(r, "contact"))
}

// Check whether Login
// If login, return user else nil
// pageType is used to navigator activated.
func checkLoginUser(r *http.Request, pageType string) interface{} {
	type tempInfo struct {
		User *User
		PageType string
	}
	if username, _ := r.Cookie("user"); username != nil {
		if users := SelectUser(map[string]interface{}{"username": username.Value}); len(users) == 0{
			return tempInfo{nil, pageType}
		} else {
			return tempInfo{&users[0], pageType}
		}

	} else {
		return tempInfo{nil, pageType}
	}
}

// Check login function.
func isLogin(r *http.Request) string {
	username, err := r.Cookie("user")
	if  err != nil {
		return "anonymous"
	} else {
		return username.Value
	}
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
	t, _ := template.ParseFiles("static/views/digest.html")
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

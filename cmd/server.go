package main

import (
	"net/http"
	"html/template"
	"fmt"
	"os"
	"io"
)

func index(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		t, _ := template.ParseFiles("./views/index.html")
		t.Execute(res, nil)
	} else {

	}
}


func Start() {
	http.Handle("/", http.RedirectHandler("/index", 301));
	http.HandleFunc("/index", index)
	http.ListenAndServe("localhost:8080", nil)
}

func download(url string, pos string) (success bool){
	if res, err := http.Get(url); err != nil {
		fmt.Printf("Download Error:", err)
	} else {
		if file, err := os.Create(pos); err != nil {
			fmt.Printf("Save File Error:", err)
		} else {
			io.Copy(file, res.Body)
			res.Body.Close()
			success = true
		}
	}
	return
}



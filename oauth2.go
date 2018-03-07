package Posger

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	oauth2Infos map[string]oauth2Info
)

func init() {
	oauth2Infos = make(map[string]oauth2Info)
	oauth2Infos["github"] = oauth2Info{"https://github.com/login/oauth/authorize", "a5676195554a7d261ec6", "http://localhost:8080/oatuh2/github/token"}
}

func registerOauth2App(router *mux.Router) {
	router.HandleFunc("/{incName}/redirect", oauth2Factory).Methods("GET")
	router.HandleFunc("/{incName}/token", oauth2FactoryToken).Methods("GET")
}

func oauth2Factory(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oauth2Url(mux.Vars(r)["incName"]), http.StatusFound)
}

func oauth2FactoryToken(w http.ResponseWriter, r *http.Request) {
	url = "https://github.com/login/oauth/access_token"
}


func oauth2Url(incName string) string {
	incInfo := oauth2Infos[incName]
	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&state=%s", incInfo.tokenUrl, incInfo.clientId, incInfo.redirectUrl, incName)
}

type oauth2Info struct {
	tokenUrl	string
	clientId 	string
	redirectUrl	string
}
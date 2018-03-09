package Posger

import (
	"strings"
	"net/url"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"io/ioutil"
)

var (
	oauth2Infos map[string]oauth2Info
)

func init() {
	oauth2Infos = make(map[string]oauth2Info)
	oauth2Infos["github"] = oauth2Info{"https://github.com/login/oauth/authorize", "https://github.com/login/oauth/access_token",
		"a5676195554a7d261ec6", "ba3d6e931b73785441aa4e1a0ab1966bc689e936", "http://localhost:8080/oauth2/github/token"}
}

func registerOauth2App(router *mux.Router) {
	router.HandleFunc("/{incName}/redirect", oauth2Factory).Methods("GET")
	router.HandleFunc("/{incName}/token", oauth2FactoryToken).Methods("GET")
}

func oauth2Factory(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oauth2Url(mux.Vars(r)["incName"]), http.StatusFound)
}

// oauth2FactoryToken is used to handle all inc.'s oauth2 code, 
// if users receive the previlege requirement, we will get token and save to mongodb.
// if other unexpected events hanpped, the paramenter's 'error' will be not null.
func oauth2FactoryToken(w http.ResponseWriter, r *http.Request) {
	paras, incName := r.URL.Query(), mux.Vars(r)["incName"]
	// if error parameter exists.
	if _, ok := paras["error"]; ok {
		http.Redirect(w, r, "/index", http.StatusFound)
		return
	}
	form := url.Values{}
	form.Add("client_id", oauth2Infos[incName].clientId)
	form.Add("client_secret", oauth2Infos[incName].clientSecret)
	form.Add("code", paras["code"][0])
	form.Add("redirect_uri", oauth2Infos[incName].redirectUrl)	// the redirectUrl should be my host index
	form.Add("state", incName)

	res, _ := http.Post(oauth2Infos[incName].tokenUrl, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Print(jsonProcessString(string(body))["access_token"])
}

// jsonProcessString segment the key-value pair to the map struct
// the format likes :
// access_token=8287e5a38e5faa9aca0a07bc522025a32cdb5b8d&scope=&token_type=bearer
func jsonProcessString(content string) (map[string]string){
	res := make(map[string]string)
	for _, segs := range strings.Split(content, "&") {
		cs := strings.Split(segs, "=")
		res[cs[0]] = cs[1]
	}
	return res
}


func oauth2Url(incName string) string {
	incInfo := oauth2Infos[incName]
	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&state=%s", incInfo.codeUrl, incInfo.clientId, incInfo.redirectUrl, incName)
}

type oauth2Info struct {
	codeUrl			string
	tokenUrl 		string
	clientId 		string
	clientSecret 	string
	redirectUrl		string
}

type oauth2Token struct {
	access_token	string
	scope 			string
	token_type		string
}
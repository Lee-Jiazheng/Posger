package Posger

import (
	"strings"
	"net/url"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"github.com/satori/go.uuid"
)

var (
	oauth2Infos map[string]oauth2Info
)

func init() {
	oauth2Infos = make(map[string]oauth2Info)
	oauth2Infos["github"] = oauth2Info{"https://github.com/login/oauth/authorize", "https://github.com/login/oauth/access_token", "https://api.github.com/user",
		"a5676195554a7d261ec6", "ba3d6e931b73785441aa4e1a0ab1966bc689e936", "http://localhost:8080/oauth2/github/token",}
}

func registeOauth2App(router *mux.Router) {
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
	res, _ = http.Get(oauth2Infos[incName].infoUrl + "?access_token=" + jsonProcessString(string(body))["access_token"])
	body, _ = ioutil.ReadAll(res.Body)
	// Get the access_token and put user information to mydatabase
	infos := &githubUser{}
	json.Unmarshal(body, &infos)
	if users := SelectUser(map[string]interface{}{"username": infos.Login, "source": incName}); len(users) == 0 {
		go AddUser(User{Source: incName, UserId: uuid.Must(uuid.NewV4()).String(), Username: infos.Login, Password: infos.Login, Avatar: infos.AvatarURL, InfoURL: infos.URL, Bio: infos.Bio})
	}

	// Later, we will marsh a better user info cookie.
	http.SetCookie(w, &http.Cookie{
		Name: "user",
		Value: infos.Login,	// user struct json
		Path: "/",
		Expires: time.Now().AddDate(0, 1, 0),
		MaxAge: 86400,	// 100 hours' validate time
	})
	http.Redirect(w, r, "/index", http.StatusFound)		// redirect to the index page
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

// The first step, redirect to the oauth2 server website.
func oauth2Url(incName string) string {
	incInfo := oauth2Infos[incName]
	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&state=%s", incInfo.codeUrl, incInfo.clientId, incInfo.redirectUrl, incName)
}

type oauth2Info struct {
	// get code's url
	codeUrl			string
	// push code to get token's url
	tokenUrl 		string
	// push access_code to get user's info's url
	infoUrl			string

	clientId 		string
	clientSecret 	string
	redirectUrl		string
}

type oauth2Token struct {
	access_token	string
	scope 			string
	token_type		string
}

type githubUser struct {
	// Following is tags.
	Login             string `json:"login"`
	ID                int    `json:"id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
	Name              string `json:"name"`
	Bio				  string `json: "bio"`
}
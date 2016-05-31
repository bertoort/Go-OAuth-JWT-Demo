package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// attributes to render a page
type attributes struct {
	Title string
	Name  string
}

// authentication json
type authentication struct {
	Authorized bool   `json:"authorized"`
	Error      string `json:"error"`
}

// index route
func index(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	attr := attributes{Title: "Index"}
	renderTemplate(res, "index", &attr)
}

// auth route for oauth authentication
func auth(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	result, err := authenticate(req)
	data := authentication{result, ""}
	if err != nil {
		data.Error = err.Error()
	}
	sendJSON(res, data)
}

// login route to store token and redirect to index
func login(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	attr := attributes{Title: "Login"}
	renderTemplate(res, "login", &attr)
}

// OAuth google login route
func oauth(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	url := getAuthURL()
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

// OAuthCallback google login route
func oauthCallback(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	code := req.URL.Query().Get("code")
	userToken := fetchToken(code)
	user := fetchUser(userToken)
	token := createJWT(user, userToken)
	url := "/login?token=" + token
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

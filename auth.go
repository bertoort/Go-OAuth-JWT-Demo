package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

// Token from oauth json
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
}

// User information from oauth
type User struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Link    string `json:"link"`
	Picture string `json:"picture"`
}

// Gplus oauth json
type Gplus struct {
	ClientID                string   `json:"client_id"`
	ProjectID               string   `json:"project_id"`
	AuthURI                 string   `json:"auth_uri"`
	TokenURI                string   `json:"token_uri"`
	ClientSecret            string   `json:"client_secret"`
	AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string   `json:"clien_x509_cert_url"`
	Userinfo                string   `json:"userinfo"`
	RedirectURIS            []string `json:"redirect_uris"`
	config                  *oauth2.Config
	prompt                  oauth2.AuthCodeOption
}

// START OAUTH CONFIG

// configures the oauth provider
func (g *Gplus) setConfig() {
	g.config = &oauth2.Config{
		ClientID:     g.ClientID,
		ClientSecret: g.ClientSecret,
		RedirectURL:  g.RedirectURIS[0],
		Endpoint: oauth2.Endpoint{
			AuthURL:  g.AuthURI,
			TokenURL: g.TokenURI,
		},
		Scopes: []string{},
	}
	g.config.Scopes = []string{"profile", "email", "openid"}
}

// END OAUTH CONFIG

var gplus Gplus
var publicKey []byte

// initializes the gplus provider
func init() {
	file, err := ioutil.ReadFile("./gplus.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(file, &gplus)
	gplus.setConfig()
	publicKey, err = ioutil.ReadFile("./rsa_key")
	if err != nil {
		log.Fatal(err)
	}
}

// gets oauth url
func getAuthURL() string {
	var opts []oauth2.AuthCodeOption
	if gplus.prompt != nil {
		opts = append(opts, gplus.prompt)
	}
	url := gplus.config.AuthCodeURL("state", opts...)
	return url
}

// START FETCH USER

// gets user information from token
func fetchUser(token Token) User {
	googleURL := gplus.Userinfo + "?access_token=" + url.QueryEscape(token.AccessToken)
	response, err := http.Get(googleURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	bits, _ := ioutil.ReadAll(response.Body)
	var user User
	json.NewDecoder(bytes.NewReader(bits)).Decode(&user)
	return user
}

// END FETCH USER

// START FETCH TOKEN
// gets token from authentication code
func fetchToken(code string) Token {
	params := url.Values{}
	params.Set("code", code)
	params.Add("client_id", gplus.ClientID)
	params.Add("client_secret", gplus.ClientSecret)
	params.Add("redirect_uri", gplus.RedirectURIS[0])
	params.Add("grant_type", "authorization_code")
	client := &http.Client{}
	req, _ := http.NewRequest("POST", gplus.TokenURI, bytes.NewBufferString(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var token Token
	json.Unmarshal(body, &token)
	return token
}

// END FETCH TOKEN

// START JWT
// creates jwt string from user and token
func createJWT(user User, userToken Token) string {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["user"] = user
	token.Claims["token"] = userToken.AccessToken
	token.Claims["exp"] = time.Now().Add(time.Second * time.Duration(userToken.ExpiresIn)).Unix()
	tokenString, err := token.SignedString(publicKey)
	if err != nil {
		log.Fatal(err)
	}
	return tokenString
}

// END JWT

// START JWT AUTH

// authenticates user token
func authenticate(req *http.Request) (bool, error) {
	token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err == nil && token.Valid {
		return true, nil
	}
	return false, err
}

// END JWT AUTH

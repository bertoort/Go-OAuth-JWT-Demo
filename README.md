# Go Oauth/JWT Demo

![jwt](https://jwt.io/img/pic_logo.svg)

Simple demo implementing [Google Plus OAuth](https://developers.google.com/identity/protocols/OAuth2) from scratch using:

- "golang.org/x/oauth2"
-	"github.com/dgrijalva/jwt-go"
- "github.com/julienschmidt/httprouter"

### Usage

- `go get github.com/bertoort/go-oauth-jwt-demo`

- Get credentials from [Google](https://console.developers.google.com/apis/library)
and add them to a `gplus.json` file

- Create a Token RSA Secret: `ssh-keygen -t rsa -f id_rsa`

- `go run *.go`

### Deployed Example

[Heroku](https://go-oauth-jwt.herokuapp.com/)

### Slides

[go-talks](http://go-talks.appspot.com/github.com/Berto/Go-OAuth-JWT-Demo/oauth-jwt.slide#1)

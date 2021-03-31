package google

import (
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
)

var (
	googleOauthConf = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/GoogleCallback",
		ClientID:     os.Getenv("GOOGLE_KEY"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		Scopes:       []string{},
		Endpoint:     google.Endpoint,
	}
	oauthRandom = "random"
)

const htmlIndex = `<html><body>
<a href="/GoogleLogin">Log in with Google</a>
</body></html>
`

func CreateClient() {
	http.HandlerFunc("/", handleMain)
	http.HandlerFunc("/GoogleLogin", handleLogin)
	http.HandlerFunc("/GoogleCallback", handleCallBack)
	fmt.Pritnln(http.ListenAndServe(":3000", nil))
}

func handleMain(writer http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

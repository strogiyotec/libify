package google

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/user"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConf = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/GoogleCallback",
		ClientID:     os.Getenv("GOOGLE_KEY"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/calendar"},
		Endpoint:     google.Endpoint,
	}
	oauthRandom = "random"
	server      http.Server
)

const htmlIndex = `<html><body>
<a href="/GoogleLogin">Log in with Google</a>
</body></html>
`

//TODO: create http client only when token doesn't exist
func CreateClient() {
	user, _ := user.Current()
	fmt.Println(user.HomeDir)
	tokenPath := user.HomeDir + "/.config/libify/token.json"
	if _, err := os.Stat(tokenPath); err == nil {
		calendarCommands := NewCalendar(tokenPath)
		calendarCommands.ListEvents()
	} else {
		//init http server
		http.HandleFunc("/", handleMain)
		http.HandleFunc("/GoogleLogin", handleLogin)
		http.HandleFunc("/GoogleCallback", handleCallBack)
		server = http.Server{Addr: ":3000", Handler: nil}
		server.ListenAndServe()
		fmt.Println(http.ListenAndServe(":3000", nil))
	}
}
func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

func handleCallBack(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthRandom {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthRandom, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	//TODO: cache this token
	token, err := googleOauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Println(token.TokenType, token.AccessToken, token.RefreshToken, token.Expiry)
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	calendar := NewCalendarFromClient(client)
	calendar.ListEvents()
	server.Shutdown(context.Background())
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConf.AuthCodeURL(oauthRandom)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

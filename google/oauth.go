package google

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
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
	if _, err := os.Stat(user.HomeDir + "/.config/libify/token.txt"); err == nil {
		fmt.Println("Exists")
	} else {
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
	//server.Shutdown()
	fmt.Println(token)
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	service, err := calendar.New(client)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	events, err := service.
		Events.
		List("primary").
		TimeMin(time.Now().
			Format(time.RFC3339)).
		MaxResults(5).
		Do()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	if len(events.Items) > 0 {
		for _, i := range events.Items {
			fmt.Fprintln(w, i.Summary, " ", i.Start.DateTime)
		}
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConf.AuthCodeURL(oauthRandom)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

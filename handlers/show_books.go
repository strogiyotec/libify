package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type LibraryCredentials struct {
	Password string
	Username string
}

//utf8=%E2%9C%93&authenticity_token=XlpUFmQajGvZlGMMwp1lbzMEdblpQBTDFU6ETOc8ZP4%3D&name=21865044497932&user_pin=5183&remember_me=true&local=false
type LoginRequest struct {
	authenticity_token string
	name               string
	user_pin           string
	remember_me        bool
	local              bool
}

func getCsrfToken() string {
	resp, err := http.Get("https://newwestminster.bibliocommons.com/user_dashboard")
	if err != nil {
		panic("Can't reach out NPL library")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Can't parse response body")
	}
	re := regexp.MustCompile(`<meta content="(?P<TOKEN>[^"]+)" name=\"csrf-token"`)
	matches := re.FindStringSubmatch(string(body))
	return matches[re.SubexpIndex("TOKEN")]
}

func (cred *LibraryCredentials) login(token string) {
	body, err := json.Marshal(LoginRequest{
		authenticity_token: token,
		local:              false,
		remember_me:        false,
		name:               cred.Username,
		user_pin:           cred.Password,
	})
	if err != nil {
		panic("Error during marshaling")
	}
	fmt.Println(body)
}

func HandleShowBooks(cred *LibraryCredentials) {
	token := getCsrfToken()
	cred.login(token)
}

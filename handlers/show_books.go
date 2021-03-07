package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

const (
	LOGIN_URL = "https://newwestminster.bibliocommons.com/user/login?destination=https://newwestminster.bibliocommons.com/"
)

type LibraryCredentials struct {
	Password string
	Username string
}

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
	fmt.Println("Send request")
	response, err := http.PostForm(LOGIN_URL, url.Values{
		"utf8":               {"%E2%9C%93"},
		"authenticity_token": {token},
		"name":               {cred.Username},
		"user_pin":           {cred.Password},
		"remember_me":        {"false"},
		"local":              {"false"}})
	if err != nil {
		panic("Error during marshaling")
	}
	defer response.Body.Close()
	fmt.Printf("Headers %v", response.Header)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("Can't read body")
	}
	fmt.Println(string(body))
}

func HandleShowBooks(cred *LibraryCredentials) {
	fmt.Println(cred)
	token := getCsrfToken()
	fmt.Printf("Token is %s\n", token)
	cred.login(token)
}

package handlers

import (
	"bytes"
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

	values := url.Values{
		"utf8":               {"%E2%9C%93"},
		"authenticity_token": {token},
		"name":               {cred.Username},
		"user_pin":           {cred.Password},
		"remember_me":        {"false"},
		"local":              {"false"}}
	req, err := http.NewRequest("POST", LOGIN_URL, bytes.NewBufferString(values.Encode()))
	if err != nil {
		panic("Error during marshaling")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Host", "newwestminster.bibliocommons.com")
	req.Header.Set("X-CSRF-Token", token)
	req.Header.Set("Origin", "https://newwestminster.bibliocommons.com")
	req.Header.Set("Referer", "https://newwestminster.bibliocommons.com/user/login?destination=https%3A%2F%2Fnewwestminster.bibliocommons.com%2F")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Printf("Headers %v", resp.Header)
	fmt.Printf("Response code %d\n", resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
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

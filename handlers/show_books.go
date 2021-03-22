package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	LOGIN_REQUEST_URL          = "https://newwestminster.bibliocommons.com/user/login?destination=%2Fuser_dashboard"
	HOME_URL                   = "https://newwestminster.bibliocommons.com/v2/checkedout/out"
	LOGIN_PAGE_URL             = "https://newwestminster.bibliocommons.com/user_dashboard"
	CHECKOUT_BOOKS_REQUEST_URL = "https://gateway.bibliocommons.com/v2/libraries/newwestminster/checkouts?accountId=%s&size=25&status=OUT&page=1&sort=status&materialType=&locale=en-CA"
	ACCESS_TOKEN               = "bc_access_token"
	SESSION_ID                 = "session_id"
)

type LibraryCredentials struct {
	Password string
	Username string
}

type HttpCredentials struct {
	bcAccessToken string
	sessionId     string
}

func (httpCred *HttpCredentials) cookies() string {
	return ACCESS_TOKEN + "=" + httpCred.bcAccessToken + ";" + SESSION_ID + "=" + httpCred.sessionId + ";"
}

func (httpCred *HttpCredentials) sendApiRequest() {
	accountId := httpCred.getAccountId()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf(CHECKOUT_BOOKS_REQUEST_URL, accountId), nil)
	req.Header.Set("Cookie", httpCred.cookies())
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	body := string(bodyBytes)
	fmt.Println(body)
}

//TODO: library doesn't recognize the request ,something wrong with cookies
func (httpCred *HttpCredentials) getAccountId() string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", HOME_URL, nil)
	req.Header.Set("Cookie", httpCred.cookies())
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	body := string(bodyBytes)
	fmt.Println(body)
	re := regexp.MustCompile(`"id":(?P<ID>[^,]+)`)
	matches := re.FindStringSubmatch(string(body))
	return matches[re.SubexpIndex("ID")]
}

func getCsrfToken() string {
	resp, err := http.Get(LOGIN_PAGE_URL)
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

func (cred *LibraryCredentials) getAccessToken(token string) HttpCredentials {
	values := url.Values{
		"name":     {cred.Username},
		"user_pin": {cred.Password}}
	req, err := http.NewRequest("POST", LOGIN_REQUEST_URL, bytes.NewBufferString(values.Encode()))
	if err != nil {
		panic("Error during marshaling")
	}
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("X-CSRF-Token", token)
	req.Header.Set("X-RESPONSIVE-PAGE", "true")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var bcToken string
	var sessionId string
	for _, value := range resp.Header["Set-Cookie"] {
		split := strings.Split(value, ";")
		if strings.Contains(split[0], ACCESS_TOKEN) {
			bcToken = split[0][strings.Index(split[0], "=")+1:]
		}
		if strings.Contains(split[0], SESSION_ID) {
			sessionId = split[0][strings.Index(split[0], "=")+1:]
		}
	}
	return HttpCredentials{bcAccessToken: bcToken, sessionId: sessionId}
}

func HandleShowBooks(cred *LibraryCredentials) {
	fmt.Println(cred)
	token := getCsrfToken()
	fmt.Printf("Token is %s\n", token)
	httpCred := cred.getAccessToken(token)
	httpCred.sendApiRequest()
}

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

//TODO: need to change this constants because they won't work for VPL
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

type Checkouts struct {
	bibTitle string
	dueDate  string
}

type HttpCredentials struct {
	bcAccessToken string
}

func (httpCred *HttpCredentials) cookies() string {
	return ACCESS_TOKEN + "=" + httpCred.bcAccessToken + ";"
}

func (httpCred *HttpCredentials) checkouts() []Checkouts {
	checkouts := []Checkouts{}
	var result map[string]interface{}
	json.Unmarshal([]byte(httpCred.sendApiRequest()), &result)
	entities := result["entities"].(map[string]interface{})
	checkoutsJsonArray := entities["checkouts"].(map[string]interface{})

	for bookId := range checkoutsJsonArray {
		// Each value is an interface{} type, that is type asserted as a string
		checkoutJson := checkoutsJsonArray[bookId].(map[string]interface{})
		checkouts = append(
			checkouts,
			Checkouts{
				bibTitle: checkoutJson["bibTitle"].(string),
				dueDate:  checkoutJson["dueDate"].(string),
			},
		)
	}
	return checkouts
}

//send api request to get a list of checkout books as json string
func (httpCred *HttpCredentials) sendApiRequest() string {
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
	return string(bodyBytes)
}

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
	re := regexp.MustCompile(`"accounts":\[(?P<ID>[^]]+)`)
	matches := re.FindStringSubmatch(string(body))
	return matches[re.SubexpIndex("ID")]
}

//Csrf token from login page
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

//Session id
func (cred *LibraryCredentials) getAccessToken(token string) HttpCredentials {
	values := url.Values{
		"name":     {cred.Username},
		"user_pin": {cred.Password}}
	req, err := http.NewRequest(
		"POST",
		LOGIN_REQUEST_URL,
		bytes.NewBufferString(values.Encode()),
	)
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
	for _, value := range resp.Header["Set-Cookie"] {
		split := strings.Split(value, ";")
		if strings.Contains(split[0], ACCESS_TOKEN) {
			bcToken = split[0][strings.Index(split[0], "=")+1:]
			break
		}
	}
	return HttpCredentials{bcAccessToken: bcToken}
}

func HandleShowBooks(cred *LibraryCredentials) {
	token := getCsrfToken()
	httpCred := cred.getAccessToken(token)
	prettyOutput(httpCred.checkouts())
}

func prettyOutput(books []Checkouts) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Title", "Due date"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.BgGreenColor},
		tablewriter.Colors{tablewriter.BgRedColor},
	)
	for _, book := range books {
		row := []string{book.bibTitle, book.dueDate}
		table.Append(row)
	}
	table.Render()

}

package google

import (
	"context"
	"encoding/json"
	"fmt"
	"libify/handlers"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

type calendarCommands struct {
	service *calendar.Service
}

//Creates a new client from given http client
func NewCalendarFromClient(client *http.Client) calendarCommands {
	service, err := calendar.New(client)
	if err != nil {
		panic(err)
	}
	calendar := calendarCommands{service: service}
	return calendar
}

//NewCalendar creates a calendar from given token json path
func NewCalendar(tokenPath string) calendarCommands {
	token := readToken(tokenPath)
	tokenSource := oauth2.ReuseTokenSource(token, oauth2.StaticTokenSource(token))
	client := oauth2.NewClient(context.Background(), tokenSource)
	service, err := calendar.New(client)
	if err != nil {
		panic(err)
	}
	calendar := calendarCommands{service: service}
	return calendar
}

func readToken(tokenPath string) *oauth2.Token {
	configFile, err := os.Open(tokenPath)
	if err != nil {
		panic(err)
	}
	//Token struct has a expire field which is a time type, can't use
	//unmarshal directly
	token := oauth2.Token{}
	var data map[string]interface{}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&data); err != nil {
		panic(err)
	}
	token.AccessToken = data["access_token"].(string)
	token.RefreshToken = data["refresh_token"].(string)
	token.TokenType = data["token_type"].(string)
	token.Expiry = time.Now().Add(time.Duration(data["expiry"].(float64)) * time.Second)
	return &token
}

func (commands *calendarCommands) ListEvents() {
	events, err := commands.service.
		Events.
		List("primary").
		TimeMin(time.Now().
			Format(time.RFC3339)).
		MaxResults(5).
		Do()
	if err != nil {

		panic(err)
	}
	if len(events.Items) > 0 {
		for _, i := range events.Items {
			fmt.Println(i.Summary, " ", i.Start.DateTime)
		}
	}
}
func (commands *calendarCommands) SaveEvents(books []handlers.Checkouts) {

}

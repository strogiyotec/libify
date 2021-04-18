package google

import (
	"context"
	"encoding/json"
	"libify/handlers"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

type calendarCommands struct {
	service *calendar.Service
}

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
	token := oauth2.Token{}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&token); err != nil {
		panic(err)
	}
	return &token
}
func (comment *calendarCommands) SaveEvents(books []handlers.Checkouts) {

}

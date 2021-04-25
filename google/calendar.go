package google

import (
	"context"
	"encoding/json"
	"libify/handlers"
	"os"
	"time"

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
func (comment *calendarCommands) SaveEvents(books []handlers.Checkouts) {

}

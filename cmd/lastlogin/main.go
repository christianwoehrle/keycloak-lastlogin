package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type Events struct {
	Realm    *string `json:"realm,omitempty"`
	Client   *string `json:"client,omitempty"`
	DateFrom *string `json:"dateFrom,omitempty"`
	DateTo   *string `json:"dateTo,omitempty"`
	Type     *string `json:"type,omitempty"`
	First    *int    `json:"first,omitempty"`
}

type EventRepresentation struct {
	UserID  *string  `json:"userId,omitempty"`
	Time    *int64   `json:"time,omitempty"`
	Type    *string  `json:"type,omitempty"`
	Details *Details `json:"details,omitempty"`
}
type Details struct {
	Username *string `json:"username,omitempty"`
}

type User struct {
	UserID    string
	Name      string
	LastLogin string
}

type AccessToken struct {
	AccessToken *string `json:"access_token,omitempty"`
}

func main() {
	user := flag.String("user", "admin", "Username to access Keycloak")
	password := flag.String("password", "xxx", "Password to access Keycloak")
	url := flag.String("url", "https://localhost:8443/", "Keycloak-URL in the form of https://localhost:8443/")
	dateFrom := flag.String("dateFrom", time.Now().AddDate(0, 0, -1).Format("2006-01-02"), "e.g. 2021-05-10")
	flag.Parse()

	client2 := resty.New().SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	resp, err := client2.R().EnableTrace().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody("grant_type=password&password=" + *password + "&username=" + *user + "&client_id=admin-cli").
		Post(*url + "auth/realms/master/protocol/openid-connect/token")
	fmt.Println(resp, err)
	fmt.Println(*dateFrom)

	body := resp.Body()
	var at AccessToken
	json.Unmarshal([]byte(body), &at)
	fmt.Println("AccessToken: ", *at.AccessToken)
	fmt.Println("=================================================")

	distinctUsers := make(map[string]User)

	complete := false
	for i := 0; !complete; i++ {

		fmt.Printf("--> %d\n", i)
		resp2, err := client2.R().
			EnableTrace().
			SetAuthToken(*at.AccessToken).
			Get(*url + "auth/admin/realms/master/events?type=LOGIN&max=100&dateFrom=" + *dateFrom + "&first=" + strconv.Itoa(i*100))

		fmt.Println(err)
		body2 := resp2.Body()
		var events []EventRepresentation
		json.Unmarshal([]byte(body2), &events)
		fmt.Println(len(events))
		if len(events) == 0 {
			complete = true
		}
		for _, u := range events {
			_, ok := distinctUsers[*u.UserID]
			if !ok {
				tm := time.Unix(*u.Time/1000, 0)
				user := User{
					UserID:    *u.UserID,
					LastLogin: tm.Format("02/01/2006, 15:04:05"),
					Name:      *u.Details.Username,
				}
				distinctUsers[user.UserID] = user
			}
		}

	}

	for i, v := range distinctUsers {
		fmt.Println(i, v)
	}
}

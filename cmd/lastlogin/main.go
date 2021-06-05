package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

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

func (u User) String() string {
	return fmt.Sprintf("%s;%s\n", u.Name, u.LastLogin)
}

type AccessToken struct {
	AccessToken *string `json:"access_token,omitempty"`
}

const batchsize = 100

func main() {
	user := flag.String("user", "admin", "Username to access Keycloak")
	password := flag.String("password", "xxx", "Password to access Keycloak")
	url := flag.String("url", "https://localhost:8443/", "Keycloak-URL in the form of https://localhost:8443/")
	dateFrom := flag.String("dateFrom", time.Now().AddDate(0, 0, -1).Format("2006-01-02"), "e.g. 2021-05-10")
	realm := flag.String("realm", "master", "e.g. master")
	logLevel := flag.String("log", "debug", "e.g. debug/fatal")

	flag.Parse()

	log.SetLevel(log.FatalLevel)
	if *logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("user: " + *user)
	log.Debugln("password: " + *password)
	log.Debugln("url: " + *url)
	log.Debugln("dateFrom: " + *dateFrom)
	log.Debugln("realm: " + *realm)
	log.Debugln("batchsize: " + strconv.Itoa(batchsize))

	client2 := resty.New().SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	resp, err := client2.R().EnableTrace().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody("grant_type=password&password=" + *password + "&username=" + *user + "&client_id=admin-cli").
		Post(*url + "auth/realms/master/protocol/openid-connect/token")

	if err != nil {
		log.Fatalln("Error accessing Keycloak Events: " + err.Error())
		panic("Exiting ...")
	}

	body := resp.Body()
	var at AccessToken
	json.Unmarshal([]byte(body), &at)
	log.Debugln("AccessToken: ", *at.AccessToken)
	log.Debugln("=================================================")

	distinctUsers := make(map[string]User)

	complete := false
	for i := 0; !complete; i++ {

		log.Debugf("ReadingEvents %d - %d", 1+(i)*batchsize, (i+1)*batchsize)
		resp2, err := client2.R().
			EnableTrace().
			SetAuthToken(*at.AccessToken).
			Get(*url + "auth/admin/realms/" + *realm + "/events?type=LOGIN&max=" + strconv.Itoa(batchsize) + "&dateFrom=" + *dateFrom + "&first=" + strconv.Itoa(i*100))

		if err != nil {
			log.Fatalln("Error accessing Keycloak Events: " + err.Error())
			panic("Exiting ...")
		}

		body2 := resp2.Body()
		var events []EventRepresentation
		json.Unmarshal([]byte(body2), &events)
		log.Debugf(" --> Events found: %d\n", len(events))
		if len(events) < batchsize {
			complete = true
		}
		for _, u := range events {
			_, ok := distinctUsers[*u.UserID]
			if !ok {
				tm := time.Unix(*u.Time/1000, 0)
				user := User{
					UserID:    *u.UserID,
					LastLogin: tm.Format(time.RFC3339),
					Name:      *u.Details.Username,
				}
				distinctUsers[user.UserID] = user
			}
		}

	}

	for _, user := range distinctUsers {
		fmt.Println(user)
	}
}

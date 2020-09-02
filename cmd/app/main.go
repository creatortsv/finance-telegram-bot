package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/creatortsv/finance-telegram-bot/internal/app/env"
	"github.com/creatortsv/finance-telegram-bot/internal/app/services/currency/exchangeratesapi"
)

type webHookRequestBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	body := &webHookRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		fmt.Println("Could not decode request body")
		return
	}

	if err := say(body.Message.Chat.ID, body.Message.Text); err != nil {
		fmt.Println("error in sending reply:", err)
		return
	}

	fmt.Println("Reply sent")
}

func say(chatID int64, input string) error {
	var a []string
	e, _ := regexp.Compile(`(change) ([0-9]{1,})([a-z]{3,3}) ([a-z]{3,3})`)
	for i, m := range e.FindStringSubmatch(strings.ToLower(input)) {
		if i > 0 {
			a = append(a, m)
		}
	}

	c, err := exchangeratesapi.New(a[2])
	if err != nil {
		return err
	}

	p, err := strconv.ParseFloat(a[1], 64)
	if err != nil {
		return err
	}

	r, err := c.Exchange(p, a[3])
	if err != nil {
		return err
	}

	body := &sendMessageReqBody{
		ChatID: chatID,
		Text:   fmt.Sprintf("%s%s is %f%s", a[1], strings.ToUpper(a[2]), r, strings.ToUpper(a[3])),
	}

	bts, err := json.Marshal(body)
	if err != nil {
		return err
	}

	res, err := http.Post("https://api.telegram.org/bot"+env.Get("TOKEN")+"/sendMessage", "application/json", bytes.NewBuffer(bts))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func main() {
	env.InitSettings()

	if err := http.ListenAndServeTLS(":"+env.Get("PORT"), "./.ssh/url_cert.pem", "./.ssh/url_private.key", http.HandlerFunc(Handler)); err != nil {
		panic(err)
	}
}

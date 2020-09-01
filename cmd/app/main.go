package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
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

var token string
var port string

func init() {
	flag.StringVar(&token, "token", "", "Telegram bot token")
	flag.StringVar(&port, "port", "80", "Port")
	flag.Parse()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	body := &webHookRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		fmt.Println("Could not decode request body")
		return
	}

	if err := say(body.Message.Chat.ID); err != nil {
		fmt.Println("error in sending reply:", err)
		return
	}

	fmt.Println("Reply sent")
}

func say(chatID int64) error {
	body := &sendMessageReqBody{
		ChatID: chatID,
		Text:   "Hello!",
	}

	bts, err := json.Marshal(body)
	if err != nil {
		return err
	}

	res, err := http.Post("https://api.telegram.org/bot"+token+"/sendMessage", "application/json", bytes.NewBuffer(bts))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func main() {
	if err := http.ListenAndServeTLS(":"+port, "./.ssh/url_cert.pem", "./.ssh/url_private.key", http.HandlerFunc(Handler)); err != nil {
		panic(err)
	}
}

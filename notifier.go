package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Web struct {
		Msg struct {
			Num0 string `yaml:"0"`
			Num1 string `yaml:"1"`
		} `yaml:"msg"`
	} `yaml:"web"`
}

type Message struct {
	Text string `json:"text"`
}

type TelegramMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func SendMessage(message interface{}, destinationApp string) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}
	if destinationApp == "telegram" {
		url := "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_APITOKEN") + "/sendMessage"
		fmt.Println(url)
		response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			return err
		}

		defer func(body io.ReadCloser) {
			if err := body.Close(); err != nil {
				log.Println("failed to close response body")
			}
		}(response.Body)
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to send successful request. Status was %q", response.Status)
		}
		return nil
	} else if destinationApp == "slack" {
		url := os.Getenv("SLACK_WEBHOOK")
		fmt.Println(url)
		response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			return err
		}

		defer func(body io.ReadCloser) {
			if err := body.Close(); err != nil {
				log.Println("failed to close response body")
			}
		}(response.Body)
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to send successful request. Status was %q", response.Status)
		}
		return nil
	} else {
		fmt.Println("Please enter the destination app: telegram or slack")
	}
	return nil
}

func ReadConfig(arg string) string {
	var config Config
	yfile, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yfile, &config)
	if err != nil {
		panic(err)
	}

	if arg == "0" {
		return config.Web.Msg.Num0
	} else {
		return config.Web.Msg.Num1
	}
}

func main() {
	if os.Args[1] == "sendMsg" {
		root := os.Getenv("CI_STATUS")
		msgText := ReadConfig(root)
		if os.Args[2] == "telegram" {
			id, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
			if err != nil {
				panic(err)
			}
			message := TelegramMessage{id, msgText}
			SendMessage(message, os.Args[2])
		} else if os.Args[2] == "slack" {
			message := Message{msgText}
			SendMessage(message, os.Args[2])
		}
	} else {
		fmt.Println("Unknown command.")
	}
}

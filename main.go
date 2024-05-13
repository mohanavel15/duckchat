package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	ENDPOINT        = "https://duckduckgo.com/duckchat/v1"
	ENDPOINT_CHAT   = ENDPOINT + "/chat"
	ENDPOINT_STATUS = ENDPOINT + "/status"
	GPT3dot5Turbo   = "gpt-3.5-turbo-0125"
	Claude3Haiku    = "claude-3-haiku-20240307"
)

func main() {
	agent := NewAgent(GPT3dot5Turbo)

	for {
		fmt.Print("user > ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		response, err := agent.Send(line)
		if err != nil {
			fmt.Println("Unexpected Error:", err.Error())
			continue
		}

		fmt.Println("Assistant >", response)
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewMessage(role, content string) Message {
	return Message{
		Role:    role,
		Content: content,
	}
}

type Chat struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

func NewChat(model string) Chat {
	return Chat{
		Model:    model,
		Messages: []Message{},
	}
}

func (c *Chat) SetMessages(messages []Message) {
	c.Messages = messages
}

func (c *Chat) AddMessage(message Message) {
	c.Messages = append(c.Messages, message)
}

func (c *Chat) Json() ([]byte, error) {
	return json.Marshal(&c)
}

type Agent struct {
	chat Chat
	vqd  string
}

func NewAgent(model string) Agent {
	a := Agent{
		chat: NewChat(model),
		vqd:  "",
	}

	a.GetVqd()

	return a
}

func (a *Agent) Send(message string) (string, error) {
	a.chat.AddMessage(NewMessage("user", message))
	return a.MakeRequest()
}

func (a *Agent) GetVqd() error {
	r, err := http.NewRequest("GET", ENDPOINT_STATUS, nil)
	if err != nil {
		return err
	}

	r.Header.Add("x-vqd-accept", "1")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	a.vqd = resp.Header.Get("x-vqd-4")
	return nil
}

func (a *Agent) MakeRequest() (string, error) {
	if a.vqd == "" {
		err := a.GetVqd()
		if err != nil {
			return "", err
		}
	}

	buf, err := a.chat.Json()
	if err != nil {
		return "", err
	}

	r, err := http.NewRequest("POST", ENDPOINT_CHAT, bytes.NewBuffer(buf))
	if err != nil {
		return "", err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("x-vqd-4", a.vqd)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return "", err
	}

	//	fmt.Println(resp.Status)

	a.vqd = resp.Header.Get("x-vqd-4")
	defer resp.Body.Close()

	ai_message := ""

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "[DONE]" {
			if ai_message == "" {
				return "", err
			}
			break
		}

		line, _ = strings.CutSuffix(line, "\n")
		line, _ = strings.CutPrefix(line, "data: ")
		if line == "" {
			continue
		}

		var res map[string]string
		json.Unmarshal([]byte(line), &res)

		if text, ok := res["message"]; ok {
			ai_message += text
		}
	}

	a.chat.AddMessage(NewMessage("assistant", ai_message))

	return ai_message, nil
}

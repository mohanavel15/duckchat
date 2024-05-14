package duckchat

import "encoding/json"

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

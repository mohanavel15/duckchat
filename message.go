package duckchat

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

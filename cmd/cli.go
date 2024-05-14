package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mohanavel15/duckchat"
)

func main() {
	agent := duckchat.NewAgent(duckchat.GPT3dot5Turbo)

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

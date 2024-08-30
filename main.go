package main

import (
	"bufio"
	"chatbot/chatbot"
	"fmt"
	"os"
	"strings"
)

var Rag chatbot.ERNIE_Rag

func QuestionInfo(input string) {
	fmt.Println(chatbot.GetColorFmt("[Question Asked]:", chatbot.ANSI_LBlue))
	fmt.Println("\"" + input + "\"")
}

func getInvalidInput() (input string) {
	fmt.Println(chatbot.GetColorFmt("[Input Question]:", chatbot.ANSI_LBlue))
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.HasSuffix(input, chatbot.AbandonSuffix) {
			fmt.Println(chatbot.GetColorFmt("[Input Question(Last Abandoned)]:", chatbot.ANSI_Green))
			continue
		}

		break
	}

	return
}

func main() {
	// we get the keys from OS enviroment variable
	Rag.SetModel(chatbot.ModelName)

	if args := os.Args; len(args) == 2 {
		input := strings.TrimSpace(args[1])
		QuestionInfo(input)
		Rag.AskQuestion(input)
		return
	}

	round := 1
	for {
		fmt.Printf(chatbot.GetColorFmt("[ROUND] %d\n", chatbot.ANSI_Yellow), round)
		input := getInvalidInput()

		if strings.HasSuffix(input, chatbot.ExitSuffix) {
			Rag.ShowTkUsage()
			fmt.Println(chatbot.GetColorFmt("[CHATBOT QUIT]", chatbot.ANSI_Red))
			return
		}

		QuestionInfo(input)
		Rag.AskQuestion(input)
		round++
	}
}

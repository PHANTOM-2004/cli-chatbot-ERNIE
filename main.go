package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var Rag ERNIE_Rag

func getInvalidInput() (input string) {
	fmt.Println(GetColorFmt("[Input Question]:", ANSI_Blue))
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.HasSuffix(input, AbandonSuffix) {
			fmt.Println(GetColorFmt("[Input Question(Last Abandoned)]:", ANSI_Green))
			continue
		}

		break
	}

	return
}

func main() {
	// we get the keys from OS enviroment variable
	Rag.SetModel(model_name)
	Rag.SetContextLimit(context_limit)

	if args := os.Args; len(args) == 2 {
		input := strings.TrimSpace(args[1])
		questionInfo(input)
		Rag.AskQuestion(input)
		return
	}

	round := 1
	for {
		fmt.Printf(GetColorFmt("[ROUND] %d\n", ANSI_Yellow), round)
		input := getInvalidInput()

		if strings.HasSuffix(input, ExitSuffix) {
      Rag.ShowTkUsage()
			fmt.Println(GetColorFmt("[CHATBOT QUIT]", ANSI_Red))
			return
		}

		questionInfo(input)
		Rag.AskQuestion(input)
		round++
	}
}

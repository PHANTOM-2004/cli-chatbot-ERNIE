package main

import (
	"bufio"
	"chatbot/chatbot"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var Rag *chatbot.ERNIE_Rag

func QuestionInfo(input string) {
	fmt.Println(chatbot.GetColorFmt("[Question Asked]:", chatbot.ANSI_LBlue))
	fmt.Println("\"" + input + "\"")
}

func getInvalidInput() string {
	fmt.Println(chatbot.GetColorFmt("[Input Question]:", chatbot.ANSI_LBlue))
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		// handle with error in input
		if err != nil {
			// log information
			Rag.ExitRound(1, err)
			// exit the program
			os.Exit(1)
		}

		input = strings.TrimSpace(input)

		if strings.HasSuffix(input, chatbot.AbandonSuffix) {
			fmt.Println(chatbot.GetColorFmt("[Input Question(Last Abandoned)]:", chatbot.ANSI_Green))
			continue
		}

		return input
	}

	log.Fatal("should not reach here")
	return ""
}

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// Run Cleanup
		fmt.Println("")
		Rag.ExitRound(1, errors.New("SIGINT"))
		os.Exit(1)
	}()
}

func main() {
  // token-count
	args := os.Args
	if len(args) >= 2 && args[1] == "--token-count" {
		chatbot.GetTotalTokensFromLog()
		return
	}

	// we get the keys from OS enviroment variable
	Rag = chatbot.NewBot(chatbot.ModelName, chatbot.ContextLimit, 1)
	// statistic at end
	defer Rag.ExitRound(1, nil)

	if len(args) == 2 {
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
			return
		}

		QuestionInfo(input)
		Rag.AskQuestion(input)
		round++
	}
}

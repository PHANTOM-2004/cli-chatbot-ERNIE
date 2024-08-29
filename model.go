package main

import (
	"context"
	"fmt"

	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
)

type ChatHistoryType = []qianfan.ChatCompletionMessage

type ERNIE_Rag struct {
	history     ChatHistoryType
	context_len int
	chat        *qianfan.ChatCompletion
	total_tks   int
}

func (Rag *ERNIE_Rag) SetModel(name string) {
	Rag.chat = qianfan.NewChatCompletion(
		qianfan.WithModel(name),
	)
	Rag.total_tks = 0
}

func (Rag *ERNIE_Rag) SetContextLimit(max_round int) {
	if max_round&1 != 0 { // when odd
		max_round += 1
	}
	Rag.context_len = max_round
}

func (Rag *ERNIE_Rag) recordA(answer string) {
	Rag.history = append(Rag.history, qianfan.ChatCompletionAssistantMessage(answer))
	if len(Rag.history) > Rag.context_len {
		panic("Should not reach here")
	}
}

func (Rag *ERNIE_Rag) recordQ(question string) {
	if len(Rag.history) > 0 && len(Rag.history) >= Rag.context_len {
		// pop the earlier QA
		Rag.history = Rag.history[2:] // abandon earlier QA(Q and A)
	}

	Rag.history = append(Rag.history, qianfan.ChatCompletionUserMessage(question))
}

func (Rag *ERNIE_Rag) ShowTkUsage() {
	fmt.Printf(GetColorFmt("[total tokens usage]: %d\n", ANSI_Red), Rag.total_tks)
}

func (Rag *ERNIE_Rag) AskQuestion(input string) {
	// now ask
	Rag.recordQ(input)
	answer := ""

	request := qianfan.ChatCompletionRequest{
		Messages:       Rag.history,
		DisableSearch:  false,
		EnableCitation: true,
	}

	// we use stream output
	response, _ := Rag.chat.Stream(
		context.TODO(),
		&request,
	)

	fmt.Printf(GetColorFmt("[%s-Answer]:\n", ANSI_Blue), model_name)

	prompt_tks, completion_tks, total_tks := 0, 0, 0
	var search_results []qianfan.SearchResult
	for {
		r, err := response.Recv()
		if err != nil {
			panic(err)
		}

		answer += r.Result
		fmt.Print(r.Result)

		if response.IsEnd {
			// jump out
			break
		}

		prompt_tks = r.Usage.PromptTokens
		completion_tks = r.Usage.CompletionTokens
		total_tks = r.Usage.TotalTokens

		if len(search_results) == 0 {
			search_results = r.SearchInfo.SearchResults
		}
	}

	// show the answer and information
	tk_fmt := GetColorFmt(tk_info_fmt, ANSI_Green)
	fmt.Printf(tk_fmt, prompt_tks, completion_tks, total_tks)
	Rag.total_tks += total_tks

	// reference list
	fmt.Println(GetColorFmt("[reference list]:", ANSI_Green))
	for i := 0; i < len(search_results); i++ {
		cur := search_results[i]
		fmt.Printf(ref_info_fmt, cur.Index, cur.Title, cur.URL)
	}
	if len(search_results) == 0 {
		fmt.Println("No reference from Internet")
	}

	Rag.recordA(answer)
}

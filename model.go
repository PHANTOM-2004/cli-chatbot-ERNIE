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
}

func (Rag *ERNIE_Rag) SetModel(name string) {
	Rag.chat = qianfan.NewChatCompletion(
		qianfan.WithModel(name),
	)
}

func (Rag *ERNIE_Rag) SetContextLimit(max_round int) {
	Rag.context_len = max_round
}

func (Rag *ERNIE_Rag) recordA(answer string) {
	Rag.history = append(Rag.history, qianfan.ChatCompletionAssistantMessage(answer))
	if len(Rag.history) > Rag.context_len {
		// pop the earlier QA
		Rag.history = Rag.history[1:]
	}
}

func (Rag *ERNIE_Rag) recordQ(question string) {
	Rag.history = append(Rag.history, qianfan.ChatCompletionUserMessage(question))
	if len(Rag.history) > Rag.context_len {
		// pop the earlier QA
		Rag.history = Rag.history[1:]
	}
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
	tk_fmt := GetColorFmt(tk_info_fmt, ANSI_Blue)
	fmt.Printf(tk_fmt, prompt_tks, completion_tks, total_tks)

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

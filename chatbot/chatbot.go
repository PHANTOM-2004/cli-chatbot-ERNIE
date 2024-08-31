package chatbot

import (
	"context"
	"fmt"

	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
)

type (
	ChatHistoryType = []qianfan.ChatCompletionMessage
	ChatType        = *qianfan.ChatCompletion
)

type ERNIE_Rag struct {
	model_name  string
	chat        ChatType
	context_len int
	total_tks   int
	history     ChatHistoryType
}

var global_logger = NewLogger(1)

func (Rag *ERNIE_Rag) SetModel(name string) {
	Rag.model_name = name
	Rag.chat = qianfan.NewChatCompletion(
		qianfan.WithModel(name),
	)
	Rag.total_tks = 0

	max_round := context_limit
	if max_round&1 != 0 { // when odd
		max_round += 1
	}
	Rag.context_len = max_round

	global_logger.LogModelConfig(1, Rag.model_name, Rag.context_len)
}

func (Rag *ERNIE_Rag) recordA(answer string) {
	Rag.history = append(Rag.history, qianfan.ChatCompletionAssistantMessage(answer))
	if len(Rag.history) > Rag.context_len {
		panic("Should not reach here")
	}

	// logging
	global_logger.LogA(1, answer, Rag.model_name)
}

func (Rag *ERNIE_Rag) recordQ(question string) {
	if len(Rag.history) > 0 && len(Rag.history) >= Rag.context_len {
		// pop the earlier QA
		Rag.history = Rag.history[2:] // abandon earlier QA(Q and A)
	}

	Rag.history = append(Rag.history, qianfan.ChatCompletionUserMessage(question))

	// logging
	global_logger.LogQ(1, question)
}

func (Rag *ERNIE_Rag) ShowTkUsage() {
	fmt.Printf(GetColorFmt("[total tokens usage]: %d\n", ANSI_Green), Rag.total_tks)
}

func (Rag *ERNIE_Rag) ExitRound(level int, err error) {
	// show statistic first
	Rag.statistic()
	global_logger.LogExitStatus(level, err)

	if err != nil {
		quit_msg := fmt.Sprintf("[CHATBOT QUIT] %s", err.Error())
		fmt.Println(GetColorFmt(quit_msg, ANSI_Red))
	} else {
		quit_msg := "[CHATBOT QUIT] normally"
		fmt.Println(GetColorFmt(quit_msg, ANSI_Green))
	}
}

func (Rag *ERNIE_Rag) statistic() {
	Rag.ShowTkUsage()
	// logging
	global_logger.LogStatistic(1, Rag.total_tks)
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

	fmt.Printf(GetColorFmt("[%s-Answer]:\n", ANSI_LBlue), ModelName)

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
	tks_output := fmt.Sprintf(tk_info_fmt, prompt_tks, completion_tks, total_tks)
	fmt.Println("\n" + GetColorFmt(tks_output, ANSI_Green))
	Rag.total_tks += total_tks

	// reference list
	var ref_output string
	for i := 0; i < len(search_results); i++ {
		cur := search_results[i]
		ref_output += fmt.Sprintf(ref_info_fmt, cur.Index, cur.Title, cur.URL)
	}
	if len(search_results) == 0 {
		ref_output = "No reference from the Internet\n"
	}

	// output reference
	fmt.Println(GetColorFmt("[reference list]:", ANSI_Green))
	fmt.Print(ref_output)

	Rag.recordA(answer)
	// logging tokens
	global_logger.Log(*NewLogEntry(1, tks_output))
	// logging reference
	global_logger.Log(*NewLogEntry(1, "[reference list]\n"+ref_output))
}

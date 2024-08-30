package chatbot

const (
	ModelName     = "ERNIE-4.0-Turbo-8K"
	AbandonSuffix = "AGAIN"
	ExitSuffix    = "QUIT"
)

const (
	tk_info_fmt = "\n[prompt tokens]: %d, " + "[completion tokens]: %d, " +
		"[total tokens]: %d\n"
	ref_info_fmt = "[%d] %s %s\n"
)

const context_limit = 4

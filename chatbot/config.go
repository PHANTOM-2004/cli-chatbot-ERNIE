package chatbot

import "os"

const (
	ModelName     = "ERNIE-4.0-Turbo-8K"
	AbandonSuffix = "AGAIN"
	ExitSuffix    = "QUIT"
	ContextLimit  = 4
)

const (
	tk_info_fmt = "[prompt tokens]: %d, " + "[completion tokens]: %d, " +
		"[total tokens]: %d"
	ref_info_fmt = "[%d] %s %s\n"
)

const (
	LogFilePerm               = 0600
	LogFileOpenMode           = os.O_CREATE | os.O_APPEND | os.O_WRONLY
	LogFileName               = ".chatbot_history"
	LogFilePrefixEnv          = "HOME"
	LogFileFallbackPathPrefix = "/home"
)

func getFilePath(prefix string, name string) string {
	return prefix + "/" + name
}

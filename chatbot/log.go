package chatbot

import (
	"fmt"
	"log"
	"os"
)

type LogEntry struct {
	Level   int // currently unused
	Message string
	// Timestamp time.Time
}

type Logger struct {
	LevelThres int
}

func NewLogger(level int) *Logger {
	return &Logger{
		LevelThres: level,
	}
}

func NewLogEntry(level int, Message string) *LogEntry {
	return &LogEntry{
		Level:   level, // currently unused
		Message: Message,
		// Timestamp: time.Now(),
	}
}

func (l *Logger) Log(entry LogEntry) {
	// we use the log path
	if entry.Level > l.LevelThres {
		// do nothing
		return
	}

	logf_path := getFilePath(os.Getenv(LogFilePrefixEnv), LogFileName)
	f, err := os.OpenFile(logf_path, LogFileOpenMode, LogFilePerm)
	if err != nil {
		fallback_path := getFilePath(LogFileFallbackPathPrefix, LogFileName)
		logf_path = fallback_path
		log.Printf(GetColorFmt("error openning default log file: [%s]; fallback: [%s]", ANSI_Red), logf_path, fallback_path)

		// open again
		f, err = os.OpenFile(fallback_path, LogFileOpenMode, LogFilePerm)
		if err != nil {
			log.Printf(GetColorFmt("error openning fallback log file: [%s]; current entry disabled", ANSI_Red), fallback_path)
		}
	}

	// log.Printf(GetColorFmt("[Logged To File]: %s\n", ANSI_Yellow), logf_path)
	defer f.Close()

	// now log
	logwrite := log.New(f, "[Log Time]: ", log.LstdFlags)
	logwrite.Println(entry.Message)
}

func (l *Logger) LogQ(level int, question string) {
	message := "[Question Asked]:\n" + question
	entry := NewLogEntry(level, message)
	l.Log(*entry)
}

func (l *Logger) LogA(level int, answer string, model_name string) {
	message := fmt.Sprintf("[%s Answered]:\n", model_name) + answer
	entry := NewLogEntry(level, message)
	l.Log(*entry)
}

func (l *Logger) LogModelConfig(level int, model_name string, context_limit int) {
	message := "[Model]: " + model_name + ", [Context Limit]: " + string(context_limit)
	entry := NewLogEntry(level, message)
	l.Log(*entry)
}

func (l *Logger) LogStatistic(level int, token_usage int) {
	message := fmt.Sprintf("[Tokens Usage]: %d\n***CHATBOT EXIT***\n", token_usage)
	entry := NewLogEntry(level, message)
	l.Log(*entry)
}

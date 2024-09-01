package chatbot

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

type LogEntry struct {
	Level   int // currently unused
	Message string
	// Timestamp time.Time
}

type Logger struct {
	LevelThres   int
	logf_handler *os.File
}

func NewLogger(level int, logf *os.File) *Logger {
	return &Logger{
		LevelThres:   level,
		logf_handler: logf,
	}
}

func (l *Logger) Close() {
	if l.logf_handler != nil {
		l.logf_handler.Close()
	}
}

func NewLogEntry(level int, Message string) *LogEntry {
	return &LogEntry{
		Level:   level, // currently unused
		Message: Message,
		// Timestamp: time.Now(),
	}
}

func getLogFileHandler(flag int) (*os.File, error) {
	logf_path := getFilePath(os.Getenv(LogFilePrefixEnv), LogFileName)
	f, err := os.OpenFile(logf_path, flag, LogFilePerm)
	if err != nil {
		fallback_path := getFilePath(LogFileFallbackPathPrefix, LogFileName)
		logf_path = fallback_path
		log.Printf(GetColorFmt("error openning default log file: [%s]; fallback: [%s]", ANSI_Red), logf_path, fallback_path)
		log.Println(GetColorFmt(err.Error(), ANSI_Red))

		// open again
		f, err = os.OpenFile(fallback_path, LogFileOpenMode, LogFilePerm)
		if err != nil {
			error_msg := fmt.Sprintf("error openning fallback log file: [%s]", fallback_path)
			return nil, errors.New(error_msg)
		}
	}
	return f, nil
}

func (l *Logger) Log(entry LogEntry) {
	// we use the log path
	if entry.Level > l.LevelThres {
		// do nothing
		return
	}

	// log.Printf(GetColorFmt("[Logged To File]: %s\n", ANSI_Yellow), logf_path)
	f := l.logf_handler
	if f == nil {
		// do nothing
		return
	}

	// now log
	logwrite := log.New(f, "", log.LstdFlags)
	logwrite.Println(entry.Message)
	// record the error
}

func (Rag *ERNIE_Rag) LogQ(level int, question string) {
	message := "[Question Asked]:\n" + question
	entry := NewLogEntry(level, message)
	l := Rag.logger
	l.Log(*entry)
}

func (Rag *ERNIE_Rag) LogA(level int, answer string, model_name string) {
	message := fmt.Sprintf("[%s Answered]:\n", model_name) + answer
	entry := NewLogEntry(level, message)
	l := Rag.logger
	l.Log(*entry)
}

func (Rag *ERNIE_Rag) LogModelConfig(level int, model_name string, context_limit int) {
	message := "[Model]: " + model_name + ", [Context Limit]: " + strconv.Itoa(context_limit)
	entry := NewLogEntry(level, message)
	l := Rag.logger
	l.Log(*entry)
}

func (Rag *ERNIE_Rag) LogExitStatus(level int, err error) {
	message_fmt := "***CHATBOT EXIT[%s]***\n"
	var message string
	if err != nil {
		message = fmt.Sprintf(message_fmt, "error: "+err.Error())
	} else {
		message = fmt.Sprintf(message_fmt, "0")
	}
	entry := NewLogEntry(level, message)
	l := Rag.logger
	l.Log(*entry)
}

func (Rag *ERNIE_Rag) LogStatistic(level int, token_usage int) {
	message := fmt.Sprintf("[Tokens Usage]: %d", token_usage)
	entry := NewLogEntry(level, message)
	l := Rag.logger
	l.Log(*entry)
}

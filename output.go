package main

import "fmt"

const (
	ANSI_Red    = "\033[31m"
	ANSI_Green  = "\033[32m"
	ANSI_Yellow = "\033[33m"
	ANSI_Blue   = "\033[94m"
	ANSI_Reset  = "\033[0m"
)

func GetColorFmt(fmt string, ansi_fmt string) string {
	return ansi_fmt + fmt + ANSI_Reset
}

func questionInfo(input string) {
	fmt.Println(GetColorFmt("[Question Asked]:", ANSI_Blue))
	fmt.Println("\"" + input + "\"")
}


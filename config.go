package main

const model_name = "ERNIE-4.0-Turbo-8K"

const (
	tk_info_fmt = "\n[prompt tokens]: %d, " + "[completion tokens]: %d, " +
		"[total tokens]: %d\n"
	ref_info_fmt  = "[%d] %s %s\n"
	AbandonSuffix = "AGAIN"
	ExitSuffix    = "QUIT"
)

const context_limit = 4

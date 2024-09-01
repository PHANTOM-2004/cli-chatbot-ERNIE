package chatbot

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

const (
	token_usage_pattern = `(.*)(\[Tokens Usage\]:\s*)(\d+)`
	token_match_index   = 3
	token_test          = "2024/09/01 12:07:02 [Tokens Usage]: 254"
)

func GetTotalTokensFromLog() {
	re, err := regexp.Compile(token_usage_pattern)
	if err != nil {
		log.Fatal(err)
	}

	// test the regular expression valid

	match := re.MatchString(token_test)
	if !match {
		log.Fatal("invalid regular expression in token count match")
	}

	f, err := getLogFileHandler(os.O_RDONLY) // read only
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	token_cnt := 0
	for scanner.Scan() {
		line := scanner.Text()
		match = re.MatchString(line)
		if !match {
			continue
		}

		single_tk_str := re.FindStringSubmatch(line)[token_match_index]
		// fmt.Println(single_tk_str)
		cnt, err := strconv.Atoi(single_tk_str)

		if cnt > 0 {
			fmt.Println(line)
		}

		if err != nil {
			log.Fatal(err)
		}

		token_cnt += cnt
	}

	fmt.Printf("%d[%.3fk] tokens", token_cnt, float64(token_cnt)/1000)
	// match
}

package service

import (
	"bufio"
	"bytes"
	"math/rand"
	"strings"
	"time"

	"github.com/haoran-mc/action-send-wework/pkg/timeutil"
	"github.com/haoran-mc/golib/pkg/log"
)

/*
00000517~我的生日
00000620~X生日
00000621~Y生日
00000622~Z生日
20250628~交房租
00001024~程序员节

1. 年份 0000 只匹配日期
*/

func ReadFormattingText(text []byte) (ret string) {
	scanner := bufio.NewScanner(bytes.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		ss := strings.Split(line, "~")
		if len(ss) != 2 {
			log.Warn("wrong line format: " + line)
			continue
		}
		t, err := timeutil.ToTime(ss[0])
		if err != nil {
			log.Warn("wrong time format", err)
			continue
		}

		// match
		now := time.Now()
		ok := t.Month() == now.Month() && t.Day() == now.Day()
		if ok && t.Year() != 0 {
			ok = t.Year() == now.Year()
		}
		if ok {
			ret += "🍭️" + ss[1] + "\n"
		}
	}
	if err := scanner.Err(); err != nil {
		log.Error("error in reading formatting text", err)
	}
	return
}

func RandomLine(text []byte) string {
	lines := bytes.Split(text, []byte("\n"))
	nonEmpty := lines[:0]
	for _, l := range lines {
		if len(bytes.TrimSpace(l)) > 0 {
			nonEmpty = append(nonEmpty, l)
		}
	}
	if len(nonEmpty) == 0 {
		return ""
	}
	return "🍭️" + string(nonEmpty[rand.Intn(len(nonEmpty))]) + "\n"
}

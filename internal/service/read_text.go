package service

import (
	"bufio"
	"bytes"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/haoran-mc/golib/pkg/log"
	"github.com/haoran-mc/golib/pkg/timeutil"
)

/*
00000517~我的生日
00000620~X生日
00000621~Y生日
00000622~Z生日
20250628~交房租
00001024~程序员节
四月廿五~我的农历生日

1. 年份 0000 只匹配日期
*/

func ReadFormattingText(text []byte) (ret []string) {
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

		ok := false
		nowyyyyMMdd := time.Now().Format("20060102")
		if startsWithDigit(ss[0]) && nowyyyyMMdd == ss[0] { // 阳历
			ok = true
		} else { // 农历
			nowLunar := timeutil.Lunar(nowyyyyMMdd)
			if strings.HasSuffix(nowLunar, ss[0]) {
				ok = true
			}
		}
		if ok {
			ret = append(ret, ss[1])
		}
	}
	if err := scanner.Err(); err != nil {
		log.Error("error in reading formatting text", err)
	}
	return
}

func startsWithDigit(s string) bool {
	for _, r := range s {
		return unicode.IsDigit(r)
	}
	return false
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
	return string(nonEmpty[rand.Intn(len(nonEmpty))])
}

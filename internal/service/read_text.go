package service

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/haoran-mc/golib/pkg/log"
	"github.com/haoran-mc/golib/pkg/timeutil"
)

/*
00000517~我的生日
20250628~交房租
00001024~程序员节
七月十五~中元节

1. 年份 0000 只匹配日期
*/

const dateFormatStr = "20060102"

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
		now := time.Now()
		fmt.Println(now)
		// now := time.Now().Add(8 * time.Hour)

		if startsWithDigit(ss[0]) { // 阳历
			t, err := time.Parse(dateFormatStr, ss[0])
			if err != nil {
				log.Warn("wrong solar date format", err)
				continue
			}
			ok = t.Month() == now.Month() && t.Day() == now.Day()
			if ok && t.Year() != 0 {
				ok = t.Year() == now.Year()
			}
		} else { // 农历
			nowyyyyMMdd := now.Format(dateFormatStr)
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

package service

import (
	"bufio"
	"bytes"
	"strings"

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
			continue
		}
		if yes, err := timeutil.TodaySpecialDate(ss[0]); err != nil && yes {
			ret += ss[1]
		}
	}
	if err := scanner.Err(); err != nil {
		log.Error("error in reading formatting text", err)
	}
	return
}

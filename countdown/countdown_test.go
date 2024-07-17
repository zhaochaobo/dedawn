package countdown_test

import (
	"dedawn/countdown"
	"gioui.org/app"
	"testing"
)

func TestTips(t *testing.T) {
	countdown.Tips("%s 后将锁屏，尽快保存好个人数据并下机")

	app.Main()
}

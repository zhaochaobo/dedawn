package control_test

import (
	"dedawn/control"
	"dedawn/mask"
	"gioui.org/app"
	"testing"
)

func TestStart(t *testing.T) {
	control.Start(&mask.Mask{})

	app.Main()
}

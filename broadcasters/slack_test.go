package broadcasters

import (
	"testing"
)

func TestD(t *testing.T) {
	slk, err := Factory("slack", map[string]interface{}{"token": "xoxb-"})
	if err != nil {
		t.Error(err.Error())
	}
	if slk != nil {
		err := slk.AddTarget("#tmplior")
		if err != nil {
			t.Error(err.Error())
		}
		err = slk.SendMessage("INFO", "hi", "bye")
		if err != nil {
			t.Error(err.Error())
		}
	}
}

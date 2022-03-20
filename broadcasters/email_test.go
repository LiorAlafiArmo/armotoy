package broadcasters

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestEmailTest(t *testing.T) {
	b, err := os.ReadFile("../config.json")
	if err != nil {
		t.Error("basa")
	}
	cfg := map[string]map[string]interface{}{}
	json.Unmarshal(b, &cfg)

	email, er := Factory("email", cfg["integrations"]["email"])
	if er != nil {
		t.Error("basa2")
	}
	fmt.Printf(" %v\n", reflect.TypeOf(email))

	concrete, _ := email.(*EMailBroadcaster)
	concrete.RecipientsBySeverity["info"] = []EmailRecipient{{Email: "lalafi@armosec.io"}}
	err = email.SendMessage("info", "test", "test")
	if err != nil {
		t.Error(err)
	}
}

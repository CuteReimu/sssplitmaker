package translate

import (
	"encoding/json"
	"os"
	"testing"
)

func TestSplits(t *testing.T) {
	buf, err := os.ReadFile("splits.json")
	if err != nil {
		t.Fatal(err)
	}
	type splitData struct {
		Key string `json:"key"`
	}
	var splits []splitData
	err = json.Unmarshal(buf, &splits)
	if err != nil {
		t.Fatal(err)
	}
	for i, split := range splits {
		switch GetIndexByID(split.Key) {
		case -1:
			t.Log("找不到分割: " + split.Key)
			t.Fail()
		case i:
		default:
			t.Logf("分割顺序错误: %s, 期望位置: %d, 实际位置: %d", split.Key, i, GetIndexByID(split.Key))
		}
	}
}

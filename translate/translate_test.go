package translate

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestSplits(t *testing.T) {
	buf, err := os.ReadFile("splits.json")
	if err != nil {
		t.Fatal(err)
	}
	type splitData struct {
		Description string `json:"description"`
		Key         string `json:"key"`
		Tooltip     string `json:"tooltip"`
		translate   string
	}
	var splits []*splitData
	err = json.Unmarshal(buf, &splits)
	if err != nil {
		t.Fatal(err)
	}
	for _, split := range splits {
		i := GetIndexByID(split.Key)
		if i == -1 {
			t.Log("找不到分割: " + split.Key)
			t.Fail()
		} else {
			split.translate = SplitsCache[i].Description
		}
	}

	// 反过来检查是否有多余
	splitMap := make(map[string]*splitData, len(splits))
	for _, split := range splits {
		splitMap[split.Key] = split
	}
	for _, split := range SplitsCache {
		if splitMap[split.ID] == nil {
			t.Log("多余分割: " + split.ID)
			t.Fail()
		}
	}
	if t.Failed() {
		t.FailNow()
	}

	fmt.Println("| Description | 翻译 | Tooltip | Key |")
	fmt.Println("|-----|------|-------------|---------|")
	for _, s := range SplitsCache {
		split := splitMap[s.ID]
		fmt.Printf("| %s | %s | %s | %s |\n", split.Description, split.translate, split.Tooltip, split.Key)
	}
}

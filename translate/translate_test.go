package translate

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSplits(t *testing.T) {
	type splitData struct {
		Alias       any    `json:"alias"`
		Description string `json:"description"`
		Key         string `json:"key"`
		Tooltip     string `json:"tooltip"`
		translate   string
	}
	var splits []*splitData
	var splitsWithAlias []*splitData
	err := json.Unmarshal(splitsJson, &splits)
	if err != nil {
		t.Fatal(err)
	}
	var newSplits []string
	for _, split := range splits {
		i := GetIndexByID(split.Key)
		if i == -1 {
			newSplits = append(newSplits, split.Key)
			t.Log("找不到分割: " + split.Key)
			t.Fail()
		} else {
			split.translate = SplitsCache[i].Description
		}
		if split.Alias != nil {
			splitsWithAlias = append(splitsWithAlias, split)
		}
	}

	// 反过来检查是否有多余
	splitMap := make(map[string]*splitData, len(splits))
	for _, split := range splits {
		splitMap[split.Key] = split
	}
	var emptyCount int
	for i, split := range SplitsCache {
		if index := GetIndexByID(split.ID); index != i && index != -1 {
			t.Log("重复分割: " + split.ID)
			t.Fail()
		}
		if split.Description == "" {
			emptyCount++
		} else if id := GetIDByDescription(split.Description); id != "" && id != split.ID {
			t.Log("重复翻译: " + split.ID)
			t.Fail()
		}
		if splitMap[split.ID] == nil {
			t.Log("多余分割: " + split.ID)
			t.Fail()
		}
	}
	t.Log("尚未翻译的数量: ", emptyCount)
	if t.Failed() {
		for _, s := range newSplits {
			fmt.Printf("{ID: \"%s\", Description: \"\"},\n", s)
		}
		t.FailNow()
	}

	fmt.Println("var cacheAliases = map[string]string{")
	for _, s := range splitsWithAlias {
		fmt.Printf("\t\"%s\": \"%s\",\n", s.Alias, s.Key)
	}
	fmt.Println("}")
}

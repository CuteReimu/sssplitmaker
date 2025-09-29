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
		Alias       any    `json:"alias"`
		Description string `json:"description"`
		Key         string `json:"key"`
		Tooltip     string `json:"tooltip"`
		translate   string
	}
	var splits []*splitData
	var splitsWithAlias []*splitData
	err = json.Unmarshal(buf, &splits)
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
	for i, split := range SplitsCache {
		if index := GetIndexByID(split.ID); index != i && index != -1 {
			t.Log("重复分割: " + split.ID)
			t.Fail()
		}
		if splitMap[split.ID] == nil {
			t.Log("多余分割: " + split.ID)
			t.Fail()
		}
	}
	if t.Failed() {
		for _, s := range newSplits {
			fmt.Printf("{\"%s\", \"\"},\n", s)
		}
		t.FailNow()
	}

	fmt.Println("var cacheAliases = map[string]string{")
	for _, s := range splitsWithAlias {
		fmt.Printf("\t\"%s\": \"%s\",\n", s.Alias, s.Key)
	}
	fmt.Println("}")

	fmt.Println("| Description | 翻译 | Tooltip | Key1 | Key2 |")
	fmt.Println("|---|---|---|---|---|")
	for _, s := range SplitsCache {
		split := splitMap[s.ID]
		var alias string
		if split.Alias != nil {
			alias = split.Alias.(string)
		}
		fmt.Printf("| %s | %s | %s | %s | %s |\n", split.Description, split.translate, split.Tooltip, split.Key, alias)
	}
}

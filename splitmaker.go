package main

import (
	"embed"
	"encoding/json"
)

//go:embed splitmaker/*
var fs embed.FS

type CategoryDirectoryData struct {
	FileName    string `json:"fileName"`
	DisplayName string `json:"displayName"`
}

func GetAllFiles() (allFiles []Option) {
	file, _ := fs.ReadFile("splitmaker/category-directory.json")
	var v map[string][]*CategoryDirectoryData
	if err := json.Unmarshal(file, &v); err != nil {
		panic(err)
	}
	for _, categoryName := range []string{"Main", "Individual Level", "Category Extensions"} {
		for _, f := range v[categoryName] {
			allFiles = append(allFiles, Option{
				Value: f.FileName + ".json",
				Label: f.DisplayName,
			})
		}
	}
	return
}

func GetSplitIds(fileName string) (string, []string, error) {
	buf, err := fs.ReadFile("splitmaker/" + fileName)
	if err != nil {
		return "", nil, err
	}
	var result struct {
		CategoryName string   `json:"categoryName"`
		Ids          []string `json:"splitIds"`
	}
	_ = json.Unmarshal(buf, &result)
	return result.CategoryName, result.Ids, nil
}

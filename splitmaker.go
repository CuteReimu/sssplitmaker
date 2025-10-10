package main

import (
	"embed"
	"encoding/json"
)

//go:embed splitmaker/*
var fs embed.FS

func GetAllFiles() (allFiles []string) {
	dir, _ := fs.ReadDir("splitmaker")
	for _, f := range dir {
		allFiles = append(allFiles, f.Name())
	}
	return
}

func GetSplitIds(fileName string) ([]string, error) {
	buf, err := fs.ReadFile("splitmaker/" + fileName)
	if err != nil {
		return nil, err
	}
	var result struct {
		Ids []string `json:"splitIds"`
	}
	_ = json.Unmarshal(buf, &result)
	return result.Ids, nil
}

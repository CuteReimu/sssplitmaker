package splitmaker

import (
	"bytes"
	"embed"
	"encoding/json"
)

//go:embed *.json
var fs embed.FS

func GetAllFiles() (allFiles []string) {
	dir, _ := fs.ReadDir(".")
	for _, f := range dir {
		allFiles = append(allFiles, f.Name())
	}
	return
}

func GetSplitIds(fileName string) ([]string, error) {
	f, err := fs.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	content := buf.String()
	var result struct {
		Ids []string `json:"splitIds"`
	}
	_ = json.Unmarshal([]byte(content), &result)
	return result.Ids, nil
}

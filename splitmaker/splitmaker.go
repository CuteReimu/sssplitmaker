package splitmaker

import (
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
	buf, err := fs.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var result struct {
		Ids []string `json:"splitIds"`
	}
	_ = json.Unmarshal(buf, &result)
	return result.Ids, nil
}

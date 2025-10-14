package cli

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

type Commit struct {
	Id        string
	Timestamp string
}

func GetLastCommit() (string, bool) {
	currentBranchName := GetCurrentBranchName()
	commits, err := os.ReadFile(".csync/branches/" + currentBranchName + "/commits.json")
	if err != nil {
		log.Fatal(err)
	}
	var content []Commit
	if err = json.Unmarshal(commits, &content); err != nil {
		log.Fatal(err)
	}
	if len(content) == 0 {
		return "", false
	}
	sort.Slice(content, func(i, j int) bool {
		return content[i].Timestamp < content[j].Timestamp
	})

	return content[len(content)-1].Id, true
}

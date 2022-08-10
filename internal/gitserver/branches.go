package gitserver

type Branch struct {
	Name    string   `json:"name"`
	IsHead  bool     `json:"isHead"`
	Commits []Commit `json:"commits"`
}

func GetBranches(username string, reponame string) ([]Branch, error) {
	var branches []Branch

	return branches, nil
}

func GetBranchByName(username string, reponame string, branchName string) (Branch, error) {
	var branch Branch

	return branch, nil
}

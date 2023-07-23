package gitdiff

import (
	gogit "github.com/go-git/go-git/v5"
)

// getFilesChangeInLastCommit will return the files changed in the last commit.
// TODO: support passing more options to the clone.
func GetChangedFilesFromLastCommit(projectPath string) ([]string, error) {

	r, err := gogit.PlainOpen(projectPath)
	if err != nil {
		return nil, err
	}

	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	// ... retrieves the commit history
	cIter, err := r.Log(&gogit.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}

	commit, err := cIter.Next()
	if err != nil {
		return nil, err
	}

	currentTree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	prevCommit, err := commit.Parent(0)
	if err != nil {
		return nil, err
	}

	prevTree, err := prevCommit.Tree()
	if err != nil {
		return nil, err
	}
	patch, err := currentTree.Patch(prevTree)
	if err != nil {
		return nil, err
	}

	var changedFiles []string
	for _, fileStat := range patch.Stats() {
		changedFiles = append(changedFiles, fileStat.Name)
	}

	return changedFiles, nil
}

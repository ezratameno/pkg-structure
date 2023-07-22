package pkgdiff

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pkgstructure "github.com/ezratameno/pkg-structure/pkg/pkg-structure"
	gogit "github.com/go-git/go-git/v5"
)

type Client struct {
	pkgStruct *pkgstructure.Client
}

func New(opts pkgstructure.Opts) (*Client, error) {
	c, err := pkgstructure.New(opts)
	if err != nil {
		return nil, err
	}

	return &Client{
		pkgStruct: c,
	}, nil
}

func (c *Client) GetPkgStructure() ([]pkgstructure.Package, error) {

	changedFiles, err := c.getFilesChangeInLastCommit()
	if err != nil {
		return nil, err
	}

	packages, err := c.pkgStruct.GetPkgStructure()
	if err != nil {
		return nil, err
	}

	pkgsChanges := make(map[string]pkgstructure.Package)

	for _, changedFile := range changedFiles {
		for _, pkg := range packages {
			for _, pkgFile := range pkg.Files {
				if strings.Contains(pkgFile, changedFile) {
					pkgsChanges[changedFile] = pkg
				}
			}
		}
	}

	// TODO: there can be multiple files from the same package that have change
	// i need to figure it out if it's a duplication or not
	for file, pkg := range pkgsChanges {
		fmt.Printf("file %s has changed in the last commit, he is a part of pkg %s\n", file, pkg.Name)
	}

	return nil, nil
}

// getFilesChangeInLastCommit will return the files changed in the last commit.
// TODO: support passing more options to the clone.
func (c *Client) getFilesChangeInLastCommit() ([]string, error) {
	r, err := gogit.PlainClone(c.pkgStruct.Module, true, &gogit.CloneOptions{
		URL: fmt.Sprintf("https://%s", c.pkgStruct.Module),
	})
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(strings.Split(c.pkgStruct.Module, string(filepath.Separator))[0])

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

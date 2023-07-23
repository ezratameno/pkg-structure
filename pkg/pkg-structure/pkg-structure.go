package pkgstructure

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Client struct {
	Opts
}

type Opts struct {
	PackagePath              string
	WithExternalDependencies bool
	Module                   string
}

func New(opts Opts) (*Client, error) {

	if opts.Module == "" {
		entries, err := os.ReadDir(opts.PackagePath)
		if err != nil {
			return nil, err
		}

		module, err := findModule(entries, opts.PackagePath)
		if err != nil {
			return nil, err
		}
		opts.Module = module
	}

	return &Client{
		Opts: opts,
	}, nil
}

func (c *Client) GetPkgStructure() ([]Package, error) {

	packages := make(map[string]Package)

	err := filepath.WalkDir(c.PackagePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}

		// ignore vendor packages
		if d.IsDir() && d.Name() == "vendor" {
			return filepath.SkipDir
		}

		// Skip non golang files
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".go" {
			return nil
		}

		// go over golang files
		file, err := c.getFileData(path)
		if err != nil {
			return err
		}

		files := []string{file.FileName}
		dependencies := file.Dependencies

		// add the information gathered from other files in of the same package
		if p, ok := packages[file.PkgName]; ok {

			files = append(files, p.Files...)
			dependencies = append(dependencies, p.Dependencies...)

			// remove duplicated dependencies between files under the same package
			dependencies = removeDuplicate[string](dependencies)
		}

		p := Package{
			Files:        files,
			Name:         file.PkgName,
			Dependencies: dependencies,
			IsMain:       file.IsMain,
		}

		packages[p.Name] = p

		return nil
	})
	if err != nil {
		return nil, err
	}

	var pkgs []Package
	for _, v := range packages {
		pkgs = append(pkgs, v)
	}

	return pkgs, nil
}

// findModule will return the module name in the go.mod file.
func findModule(items []fs.DirEntry, base string) (string, error) {
	goMod := path.Join(base, "go.mod")
	_, err := os.Stat(goMod)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("could not find go.mod file in the provided path")
		}
		return "", err
	}

	module, err := getModule(goMod)
	if err != nil {
		return "", err
	}

	return module, nil
}

func getModule(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "module") {
			return strings.TrimPrefix(line, "module "), nil
		}
	}

	return "", fmt.Errorf("module not found")
}

func (c *Client) getFileData(filePath string) (File, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return File{}, err
	}

	lines := strings.Split(string(data), "\n")

	packageName, err := getPackageName(lines)
	if err != nil {
		return File{}, err
	}

	fullPackageName := strings.TrimPrefix(filePath, c.PackagePath)
	fullPackageName = path.Dir(fullPackageName)

	// check if the directory name is not the same as the package name
	// then update the packageName
	// TODO: improve this part
	// what happens when the package is main?
	isMain := false
	if packageName == "main" {
		isMain = true
	}
	if path.Base(fullPackageName) != packageName && !isMain {
		fullPackageName = path.Join(fullPackageName, packageName)
	}

	fullPackageName = path.Join(c.Module, fullPackageName)
	p := File{
		FileName:     filePath,
		PkgName:      fullPackageName,
		Dependencies: c.getImports(lines),
		IsMain:       isMain,
	}

	return p, nil
}

func getPackageName(lines []string) (string, error) {
	for _, line := range lines {
		if strings.HasPrefix(line, "package") {
			return strings.Split(line, " ")[1], nil
		}
	}

	return "", fmt.Errorf("package not found")
}

func (c *Client) getImports(lines []string) []string {
	for i, line := range lines {
		if strings.HasPrefix(line, "import") {

			// 2 case:
			// 1. a single import
			// 2. multiple imports
			switch strings.Contains(line, `"`) {
			case true:
				// just remove the quotas around the import
				dependency := line[strings.Index(line, `"`)+1 : strings.LastIndex(line, `"`)]

				if !c.WithExternalDependencies && !strings.Contains(dependency, c.Module) {
					return nil
				}
				return []string{line[strings.Index(line, `"`)+1 : strings.LastIndex(line, `"`)]}
			case false:

				var dependencies []string
				j := i + 1
				line = lines[j]
				for !strings.Contains(line, ")") {
					line = strings.TrimSpace(line)
					if line != "" {
						dependency := line[strings.Index(line, `"`)+1 : strings.LastIndex(line, `"`)]

						if !c.WithExternalDependencies && !strings.Contains(dependency, c.Module) {

							j++
							line = lines[j]
							continue
						}

						dependencies = append(dependencies, dependency)

					}
					j++
					line = lines[j]
				}

				return dependencies
			}
		}
	}

	return nil
}

// GetDependedPackages will return depended packages.
func (c *Client) GetDependedPackages(pkg Package, pkgs []Package) []Package {
	var dependedPackages []Package

	c.getDependedPackages(pkg, pkgs, &dependedPackages)

	return dependedPackages

}

func (c *Client) getDependedPackages(pkg Package, pkgs []Package, dependedPackages *[]Package) {

	// go over all the packages and check their dependencies.
	for _, p := range pkgs {
		for _, dep := range p.Dependencies {
			if dep == pkg.Name {
				*dependedPackages = append(*dependedPackages, p)
				c.getDependedPackages(p, pkgs, dependedPackages)

			}
		}
	}

}

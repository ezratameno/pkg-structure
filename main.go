package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {

	pkgPath := flag.String("pkg-path", "/home/ezra/Desktop/golang-projects/microservices", "path to a golang project")
	flag.Parse()

	items, err := ioutil.ReadDir(*pkgPath)
	if err != nil {
		return err
	}

	module, err := findModule(items, *pkgPath)
	if err != nil {
		return err
	}

	packages := make(map[string]Package)

	err = filepath.Walk(*pkgPath, func(filePath string, info fs.FileInfo, err error) error {

		// ignore vendor packages
		if strings.Contains(filePath, "vendor") {
			return nil
		}

		// go over golang files
		if strings.HasSuffix(info.Name(), ".go") {

			pkg, err := getFileData(module, filePath)
			if err != nil {
				return err
			}

			files := []string{pkg.FileName}
			dependencies := pkg.Dependencies

			// add the information gathered from other files in of the same package
			if p, ok := packages[pkg.Name]; ok {

				files = append(files, p.Files...)
				dependencies = append(dependencies, p.Dependencies...)

				// remove duplicated dependencies between files under the same package
				dependencies = removeDuplicate[string](dependencies)
			}

			p := Package{
				Files:        files,
				Name:         pkg.Name,
				Dependencies: dependencies,
			}

			packages[p.Name] = p

		}

		return nil
	})

	if err != nil {
		return err
	}

	type wrapper struct {
		Packages map[string]Package `json:"packages"`
	}

	w := wrapper{
		Packages: packages,
	}
	b, err := json.MarshalIndent(w, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func removeDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// findModule will return the module name in the go.mod file.
func findModule(items []fs.FileInfo, base string) (string, error) {
	for _, item := range items {

		if item.IsDir() {
			continue
		}

		if item.Name() == "go.mod" {
			module, err := getModule(path.Join(base, item.Name()))
			if err != nil {
				return "", err
			}

			return module, nil
		}
	}

	return "", nil
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

type Package struct {
	Files        []string `json:"files"`
	Name         string   `json:"name"`
	Dependencies []string `json:"dependencies"`
}
type File struct {
	FileName     string   `json:"files"`
	Name         string   `json:"name"`
	Dependencies []string `json:"dependencies"`
}

func getFileData(module, filePath string) (File, error) {
	data, err := os.ReadFile(filePath)

	if err != nil {
		return File{}, err
	}

	lines := strings.Split(string(data), "\n")

	pkgName, err := getPackageName(lines)
	if err != nil {
		return File{}, err
	}

	p := File{
		FileName:     filePath,
		Name:         path.Join(module, pkgName),
		Dependencies: getImports(lines),
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

func getImports(lines []string) []string {
	for i, line := range lines {
		if strings.HasPrefix(line, "import") {

			// 2 case:
			// 1. a single import
			// 2. multiple imports
			switch strings.Contains(line, `"`) {
			case true:
				// just remove the quotas around the import

				return []string{line[strings.Index(line, `"`)+1 : strings.LastIndex(line, `"`)]}
			case false:

				var dependencies []string
				j := i + 1
				line = lines[j]
				for !strings.Contains(line, ")") {
					if strings.TrimSpace(line) != "" {
						dependencies = append(dependencies, line[strings.Index(line, `"`)+1:strings.LastIndex(line, `"`)])
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

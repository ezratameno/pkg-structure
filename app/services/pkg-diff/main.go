package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ezratameno/pkg-structure/pkg/gitdiff"
	pkgstructure "github.com/ezratameno/pkg-structure/pkg/pkg-structure"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: ", err)
		os.Exit(1)
	}
}

func run() error {
	var opts pkgstructure.Opts
	var outputType string
	var verbose bool

	flag.StringVar(&opts.PackagePath, "pkg-path", "", "Path to a golang project where the go.mod file exists. Required")
	flag.StringVar(&opts.Module, "module", "", "Full name of as appears in the go.mod file. Optional.")
	flag.BoolVar(&opts.WithExternalDependencies, "external-deps", false, "Show external dependecies to the project.")
	flag.StringVar(&outputType, "output", "plain", "Output format: <plain|json>. JSON is verbose output")
	flag.BoolVar(&verbose, "v", false, "Verbose output.")
	flag.Parse()

	if opts.PackagePath == "" {
		flag.Usage()
		return fmt.Errorf("-pkg-path is required")
	}

	if outputType != "plain" && outputType != "json" {
		return fmt.Errorf("invalid output type specified. accepeted values are plain|json")
	}

	res, err := ExtractDepsBasedOnCommitChanges(opts)
	if err != nil {
		return err
	}

	if outputType == "json" {
		return json.NewEncoder(os.Stdout).Encode(map[string]any{
			"packages": res,
		})
	}

	for _, r := range res {
		fmt.Println(r.Name)
		if verbose {
			fmt.Println("Dependencies:")
			for _, d := range r.Dependencies {
				fmt.Printf("\t%s\n", d)
			}
			fmt.Println(strings.Repeat("-", 80))

		}
	}

	return nil

}

func ExtractDepsBasedOnCommitChanges(opts pkgstructure.Opts) (map[string]pkgstructure.Package, error) {
	changedFiles, err := gitdiff.GetChangedFilesFromLastCommit(opts.PackagePath)
	if err != nil {
		return nil, err
	}

	client, err := pkgstructure.New(opts)
	if err != nil {
		return nil, err
	}

	projectPackages, err := client.GetPkgStructure()
	if err != nil {
		return nil, err
	}

	// map because a package can have multiple files that have changed
	pkgsChanges := make(map[string]pkgstructure.Package)

	// check which packages have files that have changed in the last commit
	for _, changedFile := range changedFiles {
		for _, pkg := range projectPackages {
			for _, pkgFile := range pkg.Files {
				if strings.Contains(pkgFile, changedFile) {
					pkgsChanges[changedFile] = pkg
				}
			}
		}
	}
	var pkgs []pkgstructure.Package

	// TODO: there can be multiple files from the same package that have change
	// i need to figure it out if it's a duplication or not
	for _, pkg := range pkgsChanges {
		pkgs = append(pkgs, pkg)
	}

	for _, p := range pkgs {
		dependedPackages := client.GetDependedPackages(p, projectPackages)

		if dependedPackages != nil {
			pkgs = append(pkgs, dependedPackages...)
			continue
		}
	}

	// Extract unique main packages
	mainPkgs := make(map[string]pkgstructure.Package)
	for _, p := range pkgs {
		if p.IsMain {
			mainPkgs[p.Name] = p
		}
	}

	return mainPkgs, nil
}

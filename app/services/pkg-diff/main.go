package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	pkgdiff "github.com/ezratameno/pkg-structure/internal/pkg-diff"
	pkgstructure "github.com/ezratameno/pkg-structure/pkg/pkg-structure"
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

	c := pkgdiff.New(pkgstructure.Opts{
		PackagePath:              *pkgPath,
		WithExternalDependencies: false,
	})

	s, err := c.GetPkgStructure()
	if err != nil {
		return err
	}

	type wrapper struct {
		Packages []pkgstructure.Package `json:"packages"`
	}

	w := wrapper{
		Packages: s,
	}
	b, err := json.MarshalIndent(w, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil

}

package main

import (
	"fmt"
	"os"

	pkgstructure "github.com/ezratameno/pkg-structure/pkg/pkg-structure"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {

	c, err := pkgstructure.New(pkgstructure.Opts{
		PackagePath: "/home/ezra/Desktop/golang-projects/grpc-greeter",
	})

	if err != nil {
		return err
	}

	pkgs, err := c.GetPkgStructure()
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		fmt.Printf("%+v\n", pkg)
	}
	return nil
}

package pkgstructure

type Package struct {
	Files []string `json:"files"`
	Name  string   `json:"name"`
	// Depeancies are the packages used in the 'import' of the package
	Dependencies []string `json:"dependencies"`
	IsMain       bool     `json:"isMain"`

	Types []GolangType `json:"types"`
}
type File struct {
	FileName     string   `json:"files"`
	PkgName      string   `json:"name"`
	Dependencies []string `json:"dependencies"`
	IsMain       bool     `json:"isMain"`
}

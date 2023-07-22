package pkgstructure

type Package struct {
	Files        []string `json:"files"`
	Name         string   `json:"name"`
	Dependencies []string `json:"dependencies"`
}
type File struct {
	FileName     string   `json:"files"`
	PkgName      string   `json:"name"`
	Dependencies []string `json:"dependencies"`
}

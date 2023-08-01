package pkgstructure

import (
	"os"
	"sort"
	"strings"
)

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

type GolangType struct {
	Comment string
	Name    string

	Fields []Field

	// interface/struct?
	TypeName string
}

type Field struct {
	Name    string
	Package string
	Type    string
}

func (c *Client) getTypes(pkg Package) ([]GolangType, error) {
	var pkgTypes []GolangType

	for _, f := range pkg.Files {
		content, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}

		pkgTypes = append(pkgTypes, getTypesFromFile(string(content))...)
	}

	return pkgTypes, nil
}

func getTypesFromFile(fileContent string) []GolangType {
	var pkgTypes []GolangType

	lines := strings.Split(fileContent, "\n")
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "" {
			continue
		}

		// check if it's a golang type
		// TODO: check if type is not defined in a function level
		if strings.Split(line, " ")[0] == "type" {
			typeInfo := strings.Split(line, " ")
			gType := GolangType{
				Name:     typeInfo[1],
				TypeName: typeInfo[2],
			}

			// Comment
			var commentLines []string

			// run on all the above lines to get the comments
			// we will need to reverse the lines because we went from down to up.
			for j := i - 1; j > 0; j-- {
				potentialCommentLine := lines[j]
				if strings.HasPrefix(potentialCommentLine, "//") {
					commentLine := strings.TrimPrefix(potentialCommentLine, "//")
					commentLines = append(commentLines, strings.TrimSpace(commentLine))
				} else {
					break
				}
			}

			sort.Sort(sort.Reverse(sort.StringSlice(commentLines)))
			gType.Comment = strings.Join(commentLines, "\n")

			// Fields
			var fields []Field
			for j := i; j < len(lines); j++ {
				i++

				l := strings.TrimSpace(lines[j])

				if l == "" {
					continue
				}

				//TODO: improve this part, if we have a struct inside a struct it will not work
				if l == "}" {
					break
				}

				// TODO: improve this to support interfaces and their function signature
				if gType.TypeName != "struct" {
					continue
				}

				fieldVals := strings.Split(l, " ")

				var f Field
				switch len(fieldVals) {
				// embed type
				case 1:
					// check if the type is from another package
					f.Name = fieldVals[0]
					if len(strings.Split(fieldVals[0], ".")) == 2 {
						f.Package = strings.Split(fieldVals[0], ".")[0]
						f.Type = strings.Split(fieldVals[0], ".")[1]
						f.Name = f.Type
					}

				// normal
				case 2:
					f.Name = fieldVals[0]
					f.Type = strings.Split(fieldVals[1], ".")[0]

					// check if the type is from another package
					if len(strings.Split(fieldVals[1], ".")) == 2 {
						f.Package = strings.Split(fieldVals[1], ".")[0]
						f.Type = strings.Split(fieldVals[1], ".")[1]
					}
				}

				fields = append(fields, f)

			}

			gType.Fields = fields
			pkgTypes = append(pkgTypes, gType)

		}
	}

	return pkgTypes
}

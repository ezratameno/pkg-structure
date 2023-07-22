package pkgdiff

import pkgstructure "github.com/ezratameno/pkg-structure/pkg/pkg-structure"

type Client struct {
	pkgStruct *pkgstructure.Client
}

func New(opts pkgstructure.Opts) *Client {
	return &Client{
		pkgStruct: pkgstructure.New(opts),
	}
}

func (c *Client) GetPkgStructure() ([]pkgstructure.Package, error) {
	return c.pkgStruct.GetPkgStructure()
}

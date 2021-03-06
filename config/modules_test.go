/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package config

import . "gopkg.in/check.v1"

type ModulesSuite struct{}

var _ = Suite(&ModulesSuite{})

func (s *ModulesSuite) TestValidateMissingURL(c *C) {
	m := &Submodule{Path: "foo"}
	c.Assert(m.Validate(), Equals, ErrModuleEmptyURL)
}

func (s *ModulesSuite) TestValidateBadPath(c *C) {
	input := []string{
		`..`,
		`../`,
		`../bar`,

		`/..`,
		`/../bar`,

		`foo/..`,
		`foo/../`,
		`foo/../bar`,
	}

	for _, p := range input {
		m := &Submodule{
			Path: p,
			URL:  "https://example.com/",
		}
		c.Assert(m.Validate(), Equals, ErrModuleBadPath)
	}
}

func (s *ModulesSuite) TestValidateMissingName(c *C) {
	m := &Submodule{URL: "bar"}
	c.Assert(m.Validate(), Equals, ErrModuleEmptyPath)
}

func (s *ModulesSuite) TestMarshall(c *C) {
	input := []byte(`[submodule "qux"]
	path = qux
	url = baz
	branch = bar
`)

	cfg := NewModules()
	cfg.Submodules["qux"] = &Submodule{Path: "qux", URL: "baz", Branch: "bar"}

	output, err := cfg.Marshal()
	c.Assert(err, IsNil)
	c.Assert(output, DeepEquals, input)
}

func (s *ModulesSuite) TestUnmarshall(c *C) {
	input := []byte(`[submodule "qux"]
        path = qux
        url = https://github.com/foo/qux.git
[submodule "foo/bar"]
        path = foo/bar
        url = https://github.com/foo/bar.git
		branch = dev
[submodule "suspicious"]
        path = ../../foo/bar
        url = https://github.com/foo/bar.git
`)

	cfg := NewModules()
	err := cfg.Unmarshal(input)
	c.Assert(err, IsNil)

	c.Assert(cfg.Submodules, HasLen, 2)
	c.Assert(cfg.Submodules["qux"].Name, Equals, "qux")
	c.Assert(cfg.Submodules["qux"].URL, Equals, "https://github.com/foo/qux.git")
	c.Assert(cfg.Submodules["foo/bar"].Name, Equals, "foo/bar")
	c.Assert(cfg.Submodules["foo/bar"].URL, Equals, "https://github.com/foo/bar.git")
	c.Assert(cfg.Submodules["foo/bar"].Branch, Equals, "dev")
}

func (s *ModulesSuite) TestUnmarshallMarshall(c *C) {
	input := []byte(`[submodule "foo/bar"]
	path = foo/bar
	url = https://github.com/foo/bar.git
	ignore = all
`)

	cfg := NewModules()
	err := cfg.Unmarshal(input)
	c.Assert(err, IsNil)

	output, err := cfg.Marshal()
	c.Assert(err, IsNil)
	c.Assert(string(output), DeepEquals, string(input))
}

/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package filesystem

import (
	stdioutil "io/ioutil"
	"os"

	"github.com/sniperkit/snk.fork.go-git.v4/config"
	"github.com/sniperkit/snk.fork.go-git.v4/storage/filesystem/dotgit"
	"github.com/sniperkit/snk.fork.go-git.v4/utils/ioutil"
)

type ConfigStorage struct {
	dir *dotgit.DotGit
}

func (c *ConfigStorage) Config() (conf *config.Config, err error) {
	cfg := config.NewConfig()

	f, err := c.dir.Config()
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}

		return nil, err
	}

	defer ioutil.CheckClose(f, &err)

	b, err := stdioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if err = cfg.Unmarshal(b); err != nil {
		return nil, err
	}

	return cfg, err
}

func (c *ConfigStorage) SetConfig(cfg *config.Config) (err error) {
	if err = cfg.Validate(); err != nil {
		return err
	}

	f, err := c.dir.ConfigWriter()
	if err != nil {
		return err
	}

	defer ioutil.CheckClose(f, &err)

	b, err := cfg.Marshal()
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	return err
}

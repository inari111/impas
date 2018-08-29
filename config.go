package main

import (
	"flag"
	"github.com/BurntSushi/toml"
)

var (
	// オプションを指定したときとしていないときの区別が難しいのですべて String で受ける
	optConfigFile          = flag.String("config", "./impas.toml", "config file name which includes dependency rules")
	optProjectRoot         = flag.String("root", "", `project root path from $GOROOT/src. eg. "github.com/tomoemon/impas"`)
	optIgnoreOtherProjects = flag.String("ignoreOther", "", "ignore imported packages NOT includend in the Root project if true")
	optRecursive           = flag.String("recursive", "", "search imported packages recursively if true")
)

type Constraint struct {
	From  string
	Allow []string
}

type Config struct {
	Constraint  []Constraint
	Root        string
	IgnoreOther bool
	Recursive   bool
}

func (c *Config) MaxDepth() int {
	if c.Recursive {
		return 0
	}
	return 1
}

func NewConfig() (*Config, error) {
	flag.Parse()

	c, err := LoadTOMLConfig(*optConfigFile)
	if err != nil {
		return nil, err
	}
	ApplyCommandLineOptions(c)
	return c, nil
}

func ApplyCommandLineOptions(c *Config) {
	if *optProjectRoot != "" {
		c.Root = *optProjectRoot
	}
	if *optIgnoreOtherProjects != "" {
		if *optIgnoreOtherProjects == "false" {
			c.IgnoreOther = false
		} else {
			c.IgnoreOther = true
		}
	}
	if *optRecursive != "" {
		if *optRecursive == "false" {
			c.Recursive = false
		} else {
			c.Recursive = true
		}
	}
}

func LoadTOMLConfig(fileName string) (*Config, error) {
	var config Config
	_, err := toml.DecodeFile(fileName, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
package config

import (
	"blog_agg_2/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	// Command_list []Command
	Command_map map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Command_map[name] = f

}

func (c *Commands) Run(s *State, cmd Command) error {
	err := c.Command_map[cmd.Name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

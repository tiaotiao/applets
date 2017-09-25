package commander

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type CommandFunc func(args ...string) error

type cmdMeta struct {
	Name string
	Help string
	Func CommandFunc
}

////////////////////////////////////////////////

type Commander struct {
	cmds map[string]*cmdMeta
}

func NewCommander() *Commander {
	c := &Commander{}
	c.cmds = make(map[string]*cmdMeta)
	return c
}

func (c *Commander) Register(name string, help string, f CommandFunc) error {
	name = strings.ToLower(name)
	if _, ok := c.cmds[name]; ok {
		return errors.New("Command already exist")
	}
	meta := &cmdMeta{
		Name: name,
		Help: help,
		Func: f,
	}

	c.cmds[name] = meta
	return nil
}

func (c *Commander) PrintHelp(args ...string) error {
	println(len(c.cmds), "commands:")
	for name, meta := range c.cmds {
		println("  ", name, "\t", meta.Help)
	}
	return nil
}

func (c *Commander) Run() {

	if _, ok := c.cmds["help"]; !ok {
		c.Register("help", "provide help information", c.PrintHelp)
	}
	c.Register("exit", "Quits this program.", nil)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		ss := strings.Fields(line)
		if len(ss) <= 0 {
			continue
		}
		name := strings.ToLower(ss[0])

		if name == "exit" {
			println("Bye.")
			break
		}

		args := ss[1:]
		cmd, ok := c.cmds[name]
		if !ok {
			println(name, ": Command not found. Type 'help' for more information.")
			continue
		}

		err := cmd.Func(args...)
		if err != nil {
			println("Error:", err.Error())
		}
	}
}

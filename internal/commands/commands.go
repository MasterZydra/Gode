package commands

import (
	"sort"
	"strings"
)

type Command struct {
	Name   string
	Action func()
}

type Registry struct {
	commands []Command
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) Register(command Command) {
	r.commands = append(r.commands, command)
}

func (r *Registry) Search(query string) []Command {
	query = strings.ToLower(strings.TrimSpace(query))
	results := make([]Command, 0, len(r.commands))
	for _, command := range r.commands {
		if query == "" || strings.Contains(strings.ToLower(command.Name), query) {
			results = append(results, command)
		}
	}
	sort.SliceStable(results, func(left, right int) bool {
		return results[left].Name < results[right].Name
	})
	return results
}

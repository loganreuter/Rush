package rush

import (
	"log"
	"os"
	"path"
)

type group struct {
	name   string
	path   string
	strict bool
	Schema interface{}
}

func (g *group) Member(name string) *member {
	os.MkdirAll(path.Join(g.path, name), 0755)
	if g.strict {
		return &member{name: name, schema: g.Schema, path: path.Join(g.path, name)}
	}
	return &member{name: name, path: path.Join(g.path, name)}
}

func (g *group) GetAll() []string {
	files, err := os.ReadDir(g.path)
	if err != nil {
		log.Println("\033[31m", err)
	}

	var content []string

	for _, file := range files {
		content = append(content, file.Name())
	}

	return content
}

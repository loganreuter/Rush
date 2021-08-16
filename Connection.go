package rush

import (
	"log"
	"os"
	"path"
)

type connection struct {
	dirname string
	groups  []*group
}

func Connect(relPath string, name string) (*connection, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	conn := &connection{dirname: path.Join(path.Join(dir, relPath), name)}

	if err := os.MkdirAll(conn.dirname, 0755); err != nil {
		return nil, err
	}

	return conn, nil
}

func AbsConnect(absPath string, name string) (*connection, error) {

	conn := &connection{dirname: path.Join(absPath, name)}

	if err := os.MkdirAll(conn.dirname, 0755); err != nil {
		return nil, err
	}

	return conn, nil
}

func (conn *connection) Group(name string) *group {
	for _, group := range conn.groups {
		if group.name == name {
			return group
		}
	}

	g := &group{name: name, path: path.Join(conn.dirname, name)}
	conn.groups = append(conn.groups, g)
	if err := os.MkdirAll(g.path, 0755); err != nil {
		log.Println("\033[31m", err)
	}
	return g
}

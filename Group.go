package rush

import (
	"encoding/json"
	"github.com/lreuter2020/rush/Errors"
	"log"
	"os"
	"path"
	// "reflect"
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

func (g *group) GetAll(v interface{}) {
	files, err := os.ReadDir(g.path)
	if err != nil {
		log.Println("\033[31m", err)
	}

	// if reflect.TypeOf(v).Kind() != reflect.Ptr {
	// 	error.Emit("Error: Paramater must be a pointer")
	// }
	for _, file := range files {
		log.Println(file.Type())
		if file.IsDir() {
			buff, err := os.ReadFile(path.Join(g.path, file.Name(), file.Name()+".json"))
			if err != nil {
				error.Emit(err)
			} else {
				if err := json.Unmarshal(buff, v); err != nil {
					error.Emit(err)
				}
			}

		}
		// content = append(content, file.Name())
	}

}

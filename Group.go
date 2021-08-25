package rush

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"sync"
	"time"

	CError "github.com/lreuter2020/rush/Errors"
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
	var wg sync.WaitGroup
	start := time.Now()
	files, err := os.ReadDir(g.path)
	if err != nil {
		log.Println("\033[31m", err)
	}
	var list []map[string]interface{}

	ch := make(chan map[string]interface{}, 50)

	for _, file := range files {

		if file.IsDir() {
			wg.Add(1)
			go func(file fs.DirEntry) {
				defer wg.Done()
				buff, err := os.ReadFile(path.Join(g.path, file.Name(), file.Name()+".json"))
				if err != nil {
					CError.Emit(err)
				} else {

					var rv map[string]interface{}

					if err := json.Unmarshal(buff, &rv); err != nil {
						CError.Emit(err)
					}
					ch <- rv
				}
			}(file)
		}
	}

	go func(wg *sync.WaitGroup, ch chan map[string]interface{}) {
		wg.Wait()
		close(ch)
	}(&wg, ch)

	for v := range ch {
		list = append(list, v)
	}

	blob, err := json.Marshal(list)
	if err != nil {
		CError.Emit(err)
	} else {
		json.Unmarshal(blob, v)
	}
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}

func (g *group) PipeAll(w http.ResponseWriter) {
	var wg sync.WaitGroup
	pr, pw := io.Pipe()
	pw.Write([]byte("["))
	files, err := os.ReadDir(g.path)
	if err != nil {
		CError.Emit(err)
	}

	for _, file := range files {

		if file.IsDir() {
			wg.Add(1)
			go func(file fs.DirEntry) {
				defer wg.Done()
				buff, err := os.ReadFile(path.Join(g.path, file.Name(), file.Name()+".json"))
				if err != nil {
					CError.Emit(err)
				} else {

					var rv map[string]interface{}

					if err := json.Unmarshal(buff, &rv); err != nil {
						CError.Emit(err)
					}
					json.NewEncoder(pw).Encode(&rv)
					pw.Write([]byte(","))
				}
			}(file)
		}
	}

	go func(wg *sync.WaitGroup, pw *io.PipeWriter) {
		wg.Wait()
		pw.Write([]byte("]"))
		pw.Close()
	}(&wg, pw)

	io.Copy(w, pr)
}

func (g *group) Where(args ...string) {

}

func (g *group) First(args ...interface{}) {
	if len(args) != 0 && reflect.TypeOf(args[0]).Kind() != reflect.Ptr {
		panic("First argument must be a pointer")
	}

}

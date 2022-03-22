package rush

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"

	CError "github.com/lreuter2020/rush/Errors"
)

type member struct {
	name   string
	path   string
	schema interface{}
}

func (m *member) SubGroup(name string) *group {
	if err := os.MkdirAll(path.Join(m.path, "SubGroup", name), 0777); err != nil {
		log.Println("\033[31m", err)
	}
	return &group{path: path.Join(m.path, "SubGroup", name)}
}

func (m *member) UploadFile(name string, r *http.Request, key string) {
	file, _, err := r.FormFile(key)
	if err != nil {
		panic(err)
	}

	dst, err := os.Create(path.Join(m.path, "files", name))
	if err != nil {
		panic(err)
	}

	io.Copy(dst, file)
}

func (m *member) WriteFile(name string, data []byte) {

}

func (m *member) SendFile(name string, w http.ResponseWriter) {

	file, err := os.Open(path.Join(m.path, "files", name))
	if err != nil {
		panic(err)
	}

	io.Copy(w, file)
}

func (m *member) Create(model interface{}) *member {
	// fmt.Println(reflect.TypeOf(model).String())
	if m.schema != nil {
		if !check(m.schema, model) {
			CError.Emit("Model does not match set schema")
		}
	}

	if err := os.MkdirAll(m.path, 0755); err != nil {
		CError.Emit(err)
	}

	blob, _ := json.Marshal(model)
	if err := os.WriteFile(path.Join(m.path, m.name+".json"), blob, 0666); err != nil {
		CError.Emit(err)
	}

	return m
}

func (m *member) Read(v interface{}) {
	content, err := os.ReadFile(path.Join(m.path, m.name+".json"))
	if err != nil {
		CError.Emit(err)
	}

	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		CError.Emit("Error: Must pass pointer!")
	}

	if err := json.Unmarshal(content, v); err != nil {
		CError.Emit(err)
	}
}

func (m *member) Update(model interface{}) {
	data := reflect.ValueOf(model).Interface().(map[string]interface{})

	content, err := os.ReadFile(path.Join(m.path, m.name+".json"))
	if err != nil {
		log.Println("\033[31m", err)
	}
	var h map[string]interface{}
	if err := json.Unmarshal(content, &h); err != nil {
		log.Println("\033[31m", err)
	}
	for f, v := range data {
		h[f] = v
	}

	updatedJson, _ := json.Marshal(&h)
	os.WriteFile(path.Join(m.path, m.name+".json"), updatedJson, 0777)
}

func (m *member) Destroy() {
	os.RemoveAll(m.path)
}

func check(a interface{}, b interface{}) bool {
	schema := reflect.ValueOf(a).Elem()
	h := reflect.ValueOf(b).Kind().String()
	if h == "ptr" {
		// data := reflect.ValueOf(b).Elem()
		return true
	} else {
		data := reflect.ValueOf(b).Interface().(map[string]interface{})
		// fmt.Println(data)
		for k := range data {
			if !schema.FieldByName(k).IsValid() {
				return false
			}
		}
	}

	return true
}

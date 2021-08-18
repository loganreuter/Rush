package CError

import "log"

var reset = "\033[0m"
var red = "\033[31m"

func Emit(m interface{}) {
	log.Println(red, m, reset)
}

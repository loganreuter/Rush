package error

import "log"

var reset = "\033[0m"
var red = "\033[0m"

func Emit(m interface{}) {
	log.Println(red, m, reset)
}

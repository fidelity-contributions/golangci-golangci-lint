//golangcitest:args -Egocyclo
//golangcitest:config_path testdata/gocyclo.yml
package testdata

/*
 #include <stdio.h>
 #include <stdlib.h>

 void myprint(char* s) {
 	printf("%d\n", s);
 }
*/
import "C"

import (
	"net/http"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _(s string) { // want "cyclomatic complexity .* of func .* is high .*"
	if s == http.MethodGet || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" {
		return
	}

	if s == "1" || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" {
		return
	}

	if s == "1" || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" {
		return
	}
}

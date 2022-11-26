package home

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

package lib

import (
	"fmt"
	"log"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	world := r.FormValue("world")
	log.Println("/hello endpoint world:", world)
	fmt.Fprintf(w, "Hi there, I love %s!", world)
}

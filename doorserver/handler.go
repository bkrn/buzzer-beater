//Package doorserver holds the doorbell server
package doorserver

import (
	"fmt"
	"net/http"
	"strings"
)

//Control runs the doorbell server
type Control struct {
	//http.Server
	Endpoints map[string]func(http.ResponseWriter, *http.Request)
}

//Handle accepts requests and routes them
func (cnt *Control) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ep, exists := cnt.Endpoints[strings.ToLower(r.URL.Path)]; exists {
		ep(w, r)
		return
	}
	cnt.NotFound(w, r)
	return
}

//NotFound implements 404
func (cnt *Control) NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "400 Endpoint does not exist", http.StatusBadRequest)
	return
}

//Ring accepts button presses from the AWS lambda
func Ring(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 Ring is POST only", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, "Hello!")
}

//NewControl is the control factory
func NewControl() *Control {
	return &Control{
		map[string]func(http.ResponseWriter, *http.Request){
			"/ring": Ring,
		},
	}
}

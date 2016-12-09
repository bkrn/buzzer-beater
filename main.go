//Package main launches buzzer-beater for door monitoring goodness
package main

import (
	"bkrn/buzzer-beater/doorserver"
	"net/http"
)

func main() {
	control := doorserver.NewControl()
	http.ListenAndServe(":8080", control)
}

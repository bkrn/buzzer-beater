//Package doorserver holds the doorbell server
package doorserver

import (
	"bkrn/buzzer-beater/config"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
)

//Control runs the doorbell server
type Control struct {
	//http.Server
	DB        *bolt.DB
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
func (cnt *Control) Ring(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 Ring is POST only", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, "Hello!")
}

//GetDB configures and returns the underlying Bolt DB
func GetDB() *bolt.DB {
	db, err := bolt.Open(config.Config.Database.Name, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, bucketName := range config.Config.Database.Buckets {
		db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			if err != nil {
				log.Fatal(err)
			}
			return nil
		})
	}
	return db
}

//NewControl is the control factory
func NewControl() *Control {

	cnt := &Control{
		GetDB(),
		map[string]func(http.ResponseWriter, *http.Request){},
	}

	//Declare endpoints
	cnt.Endpoints["/ring"] = cnt.Ring

	return cnt
}

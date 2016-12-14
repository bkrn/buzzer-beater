//Package doorserver holds the doorbell server
package doorserver

import (
	"bkrn/buzzer-beater/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
)

//Control runs the doorbell server
type Control struct {
	//http.Server
	DB *bolt.DB
	// Method -> URL Pattern -> handler
	Endpoints   map[string]map[string]func(http.ResponseWriter, *http.Request)
	Collections map[string]CollectionInterface
}

//Handle accepts requests and routes them
func (cnt *Control) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, exists := cnt.Endpoints[r.Method]; !exists {
		http.Error(w, fmt.Sprintf("Method %s not accepted", r.Method), http.StatusMethodNotAllowed)
		return
	}
	for pattern, ep := range cnt.Endpoints[r.Method] {
		//Error is handled in not found
		match, _ := regexp.MatchString(pattern, strings.ToLower(r.URL.String()))
		if match {
			ep(w, r)
			return
		}
	}
	cnt.NotFound(w, r)
	return
}

//RespondJSON prepares a JSON response
func (cnt *Control) RespondJSON(w http.ResponseWriter, r *http.Request, mdl interface{}) {
	jdata, err := json.MarshalIndent(mdl, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jdata)
}

//AddCollection adds an object collection
func (cnt *Control) AddCollection(name string, clc CollectionInterface) {
	cnt.Collections[name] = clc
	cnt.Endpoints["GET"][fmt.Sprintf("/%s$", name)] = func(w http.ResponseWriter, r *http.Request) {
		mdl, err := clc.All()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cnt.RespondJSON(w, r, mdl)
	}
	cnt.Endpoints["GET"][fmt.Sprintf(`/%s/(\d+)$`, name)] = func(w http.ResponseWriter, r *http.Request) {
		pattern := fmt.Sprintf(`/%s/(\d+)$`, name)
		res := regexp.MustCompile(pattern).FindStringSubmatch(strings.ToLower(r.URL.String()))[1]
		resi, err := strconv.ParseInt(res, 10, 0)
		mdl, err := clc.Get(int(resi))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cnt.RespondJSON(w, r, mdl)
	}
	cnt.Endpoints["POST"][fmt.Sprintf("/%s$", name)] = func(w http.ResponseWriter, r *http.Request) {
		bdy := []byte{}
		r.Body.Read(bdy)
		mdl, err := clc.Post(bdy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cnt.RespondJSON(w, r, mdl)
	}
	cnt.Endpoints["DELETE"][fmt.Sprintf(`/%s/(\d+)$`, name)] = func(w http.ResponseWriter, r *http.Request) {
		res := regexp.MustCompile(`/%s/(/d+)$`).FindStringSubmatch(strings.ToLower(r.URL.String()))[1]
		resi, err := strconv.ParseInt(res, 10, 0)
		mdl, err := clc.Delete(int(resi))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cnt.RespondJSON(w, r, mdl)
	}
	cnt.Endpoints["PATCH"][fmt.Sprintf(`/%s/(\d+)$`, name)] = func(w http.ResponseWriter, r *http.Request) {
		bdy := []byte{}
		r.Body.Read(bdy)
		res := regexp.MustCompile(`/%s/(/d+)$`).FindStringSubmatch(strings.ToLower(r.URL.String()))[1]
		resi, err := strconv.ParseInt(res, 10, 0)
		mdl, err := clc.Patch(int(resi), bdy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cnt.RespondJSON(w, r, mdl)
	}
}

//NotFound implements 404
func (cnt *Control) NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("400 resource does not exist at %s", strings.ToLower(r.URL.String())), http.StatusBadRequest)
	return
}

//Ring accepts button presses from the AWS lambda
func (cnt *Control) Ring(w http.ResponseWriter, r *http.Request) {
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
		map[string]map[string]func(http.ResponseWriter, *http.Request){
			"GET":    map[string]func(http.ResponseWriter, *http.Request){},
			"POST":   map[string]func(http.ResponseWriter, *http.Request){},
			"PATCH":  map[string]func(http.ResponseWriter, *http.Request){},
			"PUT":    map[string]func(http.ResponseWriter, *http.Request){},
			"DELETE": map[string]func(http.ResponseWriter, *http.Request){},
		},
		map[string]CollectionInterface{},
	}

	//Declare endpoints
	cnt.Endpoints["POST"]["/ring$"] = cnt.Ring
	cnt.AddCollection("users", &UserCollection{cnt.DB})
	//cnt.AddCollection("messages", &MessageCollection{cnt.DB})

	return cnt
}

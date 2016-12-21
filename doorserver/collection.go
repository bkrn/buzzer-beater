package doorserver

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"strings"

	"github.com/asaskevich/govalidator"
	bolt "github.com/boltdb/bolt"
)

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

//CollectionInterface is for collecitons. Implements all standard REST Endpoints
//defferentiating between GET with ID (one) and GET without (all)
type CollectionInterface interface {
	//Rest methods
	Get(int) (interface{}, error)
	All() (interface{}, error)
	Post([]byte) (interface{}, error)
	Patch(int, []byte) (interface{}, error)
	Delete(int) (interface{}, error)
	//Query methods
	//TODO GET or ALL with url parameters filters in find by field
	FindByField(string, string) (interface{}, error)
}

//DoorCollection is the collection base struct
type DoorCollection struct {
	DB             *bolt.DB
	CollectionName string
	GetModel       func() DoorModelInterface
}

//Get returns one user from ID
func (clc *DoorCollection) Get(ID int) (interface{}, error) {
	mdl := clc.GetModel()
	err := clc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(clc.CollectionName))
		err := json.Unmarshal(b.Get(itob(ID)), mdl)
		return err
	})
	return mdl, err
}

//All returns all users from db
func (clc *DoorCollection) All() (interface{}, error) {
	mdls := []interface{}{}
	err := clc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(clc.CollectionName))
		b.ForEach(func(k, v []byte) error {
			mdl := clc.GetModel()
			err := json.Unmarshal(v, mdl)
			mdls = append(mdls, mdl)
			return err
		})
		return nil
	})
	return mdls, err
}

//Post adds a new user the to database and returns it
func (clc *DoorCollection) Post(data []byte) (interface{}, error) {
	mdl := clc.GetModel()
	err := json.Unmarshal(data, mdl)
	if err != nil {
		return mdl, err
	}
	valid, err := govalidator.ValidateStruct(mdl)
	if err != nil {
		return mdl, err
	}
	if !valid {
		return mdl, errors.New("Invalid user object")
	}
	err = clc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(clc.CollectionName))

		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()
		mdl.SetID(int(id))
		var buf []byte
		buf, err = json.Marshal(mdl)
		if err != nil {
			return err
		}

		return b.Put(itob(mdl.GetID()), buf)
	})
	return mdl, err
}

//Delete returns
func (clc *DoorCollection) Delete(ID int) (interface{}, error) {
	mdl := clc.GetModel()
	err := clc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(clc.CollectionName))
		err := b.Delete(itob(ID))
		return err
	})
	return mdl, err
}

//Patch updates a user
func (clc *DoorCollection) Patch(ID int, data []byte) (interface{}, error) {
	mdl := clc.GetModel()
	err := clc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(clc.CollectionName))
		//Get origional values
		err := json.Unmarshal(b.Get(itob(ID)), mdl)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, mdl)
		if err != nil {
			return err
		}
		buf, err := json.Marshal(mdl)
		if err != nil {
			return err
		}
		return b.Put(itob(mdl.GetID()), buf)
	})
	return mdl, err
}

//FindByField finds the first model by a field
func (clc *DoorCollection) FindByField(field string, value string) (interface{}, error) {
	mdl := clc.GetModel()
	err := clc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(clc.CollectionName))
		var found bool
		b.ForEach(func(k, v []byte) error {
			if !found {
				cmap := map[string]interface{}{}
				err := json.Unmarshal(v, &cmap)
				if err != nil {
					return err
				}
				mdlVal, fieldFound := cmap[strings.ToLower(field)]
				if fieldFound && mdlVal.(string) == value {
					err := json.Unmarshal(v, mdl)
					if err != nil {
						return err
					}
					found = true
				}
			}
			return nil
		})
		if !found {
			return errors.New("Model not found")
		}
		return nil
	})
	if err != nil {
		return clc.GetModel(), err
	}
	return mdl, err
}

//NewUserCollection returns a user collection
func NewUserCollection(db *bolt.DB) *DoorCollection {
	clc := &DoorCollection{
		db, "User", func() DoorModelInterface { return &DoorUser{} },
	}
	return clc
}

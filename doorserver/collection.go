package doorserver

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/boltdb/bolt"
)

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

//CollectionInterface is for collecitons. Implements all standard REST Endpoints
//defferentiating between GET with ID (one) and GET without (all)
type CollectionInterface interface {
	Get(int) (interface{}, error)
	All() (interface{}, error)
	Post([]byte) (interface{}, error)
	Patch(int, []byte) (interface{}, error)
	Delete(int) (interface{}, error)
}

//UserCollection interacts with users
type UserCollection struct {
	DB *bolt.DB
}

//Get returns one user from ID
func (clc *UserCollection) Get(ID int) (interface{}, error) {
	mdl := &DoorUser{}
	err := clc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("User"))
		err := json.Unmarshal(b.Get(itob(ID)), mdl)
		return err
	})
	return mdl, err
}

//All returns all users from db
func (clc *UserCollection) All() (interface{}, error) {
	mdls := []interface{}{}
	err := clc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("User"))
		b.ForEach(func(k, v []byte) error {
			mdl := &DoorUser{}
			err := json.Unmarshal(v, mdl)
			mdls = append(mdls, mdl)
			return err
		})
		return nil
	})
	return mdls, err
}

//Post adds a new user the to database and returns it
func (clc *UserCollection) Post(data []byte) (interface{}, error) {
	mdl := &DoorUser{}
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
		b := tx.Bucket([]byte("User"))

		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()
		mdl.ID = int(id)
		var buf []byte
		buf, err = json.Marshal(mdl)
		if err != nil {
			return err
		}

		return b.Put(itob(mdl.ID), buf)
	})
	return mdl, err
}

//Delete returns
func (clc *UserCollection) Delete(ID int) (interface{}, error) {
	mdl := &DoorUser{}
	err := clc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("User"))
		err := b.Delete(itob(ID))
		return err
	})
	return mdl, err
}

//Patch updates a user
func (clc *UserCollection) Patch(ID int, data []byte) (interface{}, error) {
	mdl := &DoorUser{}
	err := clc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("User"))
		//Get origional values
		err := json.Unmarshal(b.Get(itob(mdl.ID)), mdl)
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
		return b.Put(itob(mdl.ID), buf)
	})
	return mdl, err
}

package doorserver

import "golang.org/x/crypto/bcrypt"

//DoorModelInterface is common model interface
type DoorModelInterface interface {
	GetID() int
	SetID(int)
}

//DoorModel is common model struct
type DoorModel struct {
	ID int `json:"id" valid:"-"`
}

//GetID gets the model ID
func (mdl *DoorModel) GetID() int {
	return mdl.ID
}

//SetID sets the model ID
func (mdl *DoorModel) SetID(nid int) {
	mdl.ID = nid
}

//DoorUser stores information about those allowed to answer the door
type DoorUser struct {
	DoorModel
	Name     string `json:"name" valid:"alphanum,required"`
	Phone    string `json:"phone,omitempty" valid:"optional"`
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required"`
}

//HashPass hashes a user's plain text password or returns an error
func (user *DoorUser) HashPass() error {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.Password = string(pass[:])
	return nil
}

//Authenticate returns nil if pw is valid
func (user *DoorUser) Authenticate(pw string) error {
	pwb := []byte(pw)
	hsh := []byte(user.Password)
	return bcrypt.CompareHashAndPassword(hsh, pwb)
}

//DoorMessage stores messages that are served to those ringing the door
type DoorMessage struct {
	DoorModel
	Name  string `json:"name" valid:"alphanum,required"`
	Image string `json:"image" valid:"optional"`
	Text  string `json:"text" valid:"alphanum,required"`
}

package doorserver

import "golang.org/x/crypto/bcrypt"

type doorModel struct {
	ID int `valid:"-"`
}

//ModelInterface

//DoorUser stores information about those allowed to answer the door
type DoorUser struct {
	doorModel
	Name     string `json:"name" valid:"alphanum,required"`
	Phone    string `json:"phone,omitempty" valid:"optional"`
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"-" valid:"required"`
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

//DoorSession holds a user session information
type DoorSession struct {
	doorModel
}

//DoorMessage hold a message that can be used by the door
type DoorMessage struct {
	doorModel
}

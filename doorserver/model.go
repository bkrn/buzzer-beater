package doorserver

//DoorUser stores information about those allowed to answer the door
type DoorUser struct {
	ID       int
	Name     string
	Phone    string
	Email    string
	Password string
}

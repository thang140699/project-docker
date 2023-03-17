package model

type UserDetails struct {
	Name     string
	Email    string
	Password string
}

// FacebookUserDetails is struct used for user details
type FacebookUserDetails struct {
	ID    string
	Name  string
	Email string
}

package model

type AuthSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AuthSignUp struct {
	IdentityID  string   `bson:"identityID" json:"identityID"`
	FirstName   string   `bson:"firstName" json:"firstName"`
	LastName    string   `bson:"lastName" json:"lastName"`
	Birthday    int      `bson:"birthday" json:"birthday"`
	Gender      int      `bson:"gender" json:"gender"`
	PhoneNumber string   `bson:"phoneNumber" json:"phoneNumber"`
	Locations   []string `bson:"locations" json:"locations"`
	Email       string   `bson:"email" json:"email" valid:"required,email"`
	Password    string   `bson:"password" json:"password" valid:"required"`
	Scopes      string   `bson:"scopes" json:"scopes"`
	UserID      string   `bson:"userID" json:"userID"`
}
type ResetPassword struct {
	Contact         string
	Password        string
	ConfirmPassword string
}
type ChangePassword struct {
	Contact         string
	Password        string
	NewPassword     string
	ConfirmPassword string
}

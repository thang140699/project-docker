package googleUltis

import (
	"WeddingBackEnd/model"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
	googleOAuth "golang.org/x/oauth2/google"
	"gopkg.in/mgo.v2/bson"
)

type GoogleConfig struct {
	Client_Id     string
	Client_Secret string
	Redirect_Url  string
}

type GoogleUser struct {
	ID             string `bson:"id" json:"id"`
	Email          string `bson:"email" json:"email"`
	Verified_email bool   `bson:"verified_email" json:"verified_email"`
	Name           string `bson:"name" json:"name"`
	Given_name     string `bson:"given_name" json:"given_name"`
	Family_name    string `bson:"family_name" json:"family_name"`
	Picture        string `bson:"picture" json:"picture"`
	Locale         string `bson:"locale" json:"locale"`
}

func GoogleOAuthConfig() *oauth2.Config {
	h := new(GoogleConfig)
	return &oauth2.Config{
		ClientID:     h.Client_Id,
		ClientSecret: h.Client_Secret,
		RedirectURL:  h.Redirect_Url,
		Endpoint:     googleOAuth.Endpoint,
		Scopes:       []string{"email"},
	}
}

func GetUserInfoFromGoogle(token string) (model.User, error) {
	var userGoogleInfo GoogleUser
	resp, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo?access_token="+token, nil)
	respUserDetail, err := http.DefaultClient.Do(resp)
	if err != nil {
		return model.User{}, errors.New("error occurred while getting information from Google")
	}

	err = json.NewDecoder(respUserDetail.Body).Decode(&userGoogleInfo)
	defer respUserDetail.Body.Close()
	if err != nil {
		return model.User{}, errors.New("error occurred while getting information from Google")
	}

	User := model.User{
		ID:          bson.ObjectId(userGoogleInfo.ID),
		IdentityID:  "",
		FirstName:   userGoogleInfo.Given_name,
		LastName:    userGoogleInfo.Family_name,
		Birthday:    0,
		Gender:      -1,
		PhoneNumber: "",
		Email:       userGoogleInfo.Email,
	}

	return User, nil

}

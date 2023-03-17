package facbookUltilities

import (
	"WeddingBackEnd/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	facebookOAuth "golang.org/x/oauth2/facebook"
	"gopkg.in/mgo.v2/bson"
)

type FacebookConfig struct {
	Client_Id     string
	Client_Secret string
	Redirect_Url  string
}

func FacebookOAuthConfig() *oauth2.Config {
	h := new(FacebookConfig)
	return &oauth2.Config{
		ClientID:     h.Client_Id,
		ClientSecret: h.Client_Secret,
		RedirectURL:  h.Redirect_Url,
		Endpoint:     facebookOAuth.Endpoint,
		Scopes:       []string{"email"},
	}
}

func GetUserInfoFromFacebook(token string) (model.User, error) {
	var fbUserDetails model.FacebookUserDetails
	fbUserDetailRequest, _ := http.NewRequest("GET", "https://graph.facebook.com/me?fields=id,name,email&access_token="+token, nil)
	facebookUserDetailsResponse, err := http.DefaultClient.Do(fbUserDetailRequest)
	if err != nil {
		fmt.Println("response : ", err)
		return model.User{}, errors.New("error occurred while getting information from Facebook")
	}
	err = json.NewDecoder(facebookUserDetailsResponse.Body).Decode(&fbUserDetails)
	defer facebookUserDetailsResponse.Body.Close()
	if err != nil {
		return model.User{}, errors.New("error decode while getting information from Facebook")
	}
	user := model.User{
		ID:          bson.NewObjectId(),
		IdentityID:  "",
		FirstName:   fbUserDetails.Name,
		LastName:    fbUserDetails.Name,
		Birthday:    0,
		Gender:      0,
		PhoneNumber: "",
		Email:       fbUserDetails.Email,
		Locations:   []model.Location{},
		CreatedTime: time.Now().UTC(),
		UpdatedTime: time.Now().UTC(),
	}
	return user, nil
}

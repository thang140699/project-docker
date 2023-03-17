package mongo

import (
	"WeddingBackEnd/model"
	"WeddingBackEnd/ultilities/provider/mongo"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const AccountMongoCollection = "account"

type AccountMongoRepository struct {
	provider       *mongo.MongoProvider
	collectionName string
}

func NewAccountMongoRepository(provider *mongo.MongoProvider) *AccountMongoRepository {
	repo := &AccountMongoRepository{provider, AccountMongoCollection}
	collection, close := repo.collection()
	defer close()

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"email",
		},
		Unique: true,
	})

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"userID",
		},
		Unique: true,
	})

	return repo
}

func (repo *AccountMongoRepository) collection() (collection *mgo.Collection, close func()) {
	session := repo.provider.MongoClient().GetCopySession()
	close = session.Close

	return session.DB(repo.provider.MongoClient().Database()).C(repo.collectionName), close
}
func (repo *AccountMongoRepository) All() ([]model.Account, error) {
	collection, close := repo.collection()
	defer close()

	result := make([]model.Account, 0)
	err := collection.Find(nil).All(&result)
	return result, repo.provider.NewError(err)
}

func (repo *AccountMongoRepository) FindByID(id string) (*model.Account, error) {
	collection, close := repo.collection()
	defer close()

	if !bson.IsObjectIdHex(id) {
		return nil, fmt.Errorf("invalid id")
	}

	var account model.Account
	err := collection.Find(bson.ObjectIdHex(id)).One(&account)
	return &account, repo.provider.NewError(err)
}

func (repo *AccountMongoRepository) FindByEmail(email string) (*model.Account, error) {
	collection, close := repo.collection()
	defer close()

	var account model.Account
	err := collection.Find(bson.M{"email": email}).One(&account)
	return &account, repo.provider.NewError(err)
}

func (repo *AccountMongoRepository) FindByPhoneNumber(phoneNumber string) (*model.Account, error) {
	collection, close := repo.collection()
	defer close()

	var account model.Account
	err := collection.Find(bson.M{"phoneNumber": phoneNumber}).One(&account)
	return &account, repo.provider.NewError(err)
}

func (repo *AccountMongoRepository) Save(account model.Account) error {
	collection, close := repo.collection()
	defer close()

	err := collection.Insert(account)
	return repo.provider.NewError(err)

}

func (repo *AccountMongoRepository) UpdateByEmail(email string, account model.Account) error {
	collection, close := repo.collection()
	defer close()

	err := collection.Update(bson.M{"email": email}, account)
	return repo.provider.NewError(err)
}

func (repo *AccountMongoRepository) UpdateByPhoneNumber(phoneNumber string, account model.Account) error {
	collection, close := repo.collection()
	defer close()

	err := collection.Update(bson.M{"phoneNumber": phoneNumber}, account)
	return repo.provider.NewError(err)
}

func (repo *AccountMongoRepository) RemoveByID(id string) error {
	collection, close := repo.collection()
	defer close()

	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("invalid id")
	}

	err := collection.RemoveId(bson.ObjectIdHex(id))
	return repo.provider.NewError(err)
}

func (repo *AccountMongoRepository) RemoveByEmail(email string) error {
	collection, close := repo.collection()
	defer close()

	err := collection.Remove(bson.M{"email": email})
	return repo.provider.NewError(err)
}

func (repo *AccountMongoRepository) RemoveByUserID(userID string) error {
	collection, close := repo.collection()
	defer close()

	err := collection.Remove(bson.M{"userID": userID})
	return repo.provider.NewError(err)
}

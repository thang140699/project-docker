package mongo

import (
	"errors"
	"log"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	DB_SERVER          = "DB_SERVER"
	DB_USERNAME        = "DB_USERNAME"
	DB_PASSWORD        = "DB_PASSWORD"
	DB_SOURCE          = "DB_SOURCE"
	DB_TIMEOUT         = "DB_TIMEOUT"
	DB_DATABASE        = "DB_DATABASE"
	DB_COLLECTION      = "DB_COLLECTION"
	DB_DEFAULT_TIMEOUT = 10000
)

type MongoDB struct {
	addrs      []string
	username   string
	password   string
	source     string
	timeout    int
	session    *mgo.Session
	database   string
	collection string
	url        string
}

var isLocked = false

func NewMongoDBFromURL(u string) *MongoDB {
	timeout := DB_DEFAULT_TIMEOUT

	instance := MongoDB{
		url:     u,
		timeout: timeout,
	}

	err := instance.Init()
	if err != nil {
		log.Println(err)
		return nil
	}

	return &instance
}

func NewMongoDB(config map[string]string) *MongoDB {
	timeout, err := strconv.Atoi(config[DB_TIMEOUT])
	if err != nil {
		timeout = DB_DEFAULT_TIMEOUT
	}

	instance := MongoDB{
		addrs:      []string{config[DB_SERVER]},
		username:   config[DB_USERNAME],
		password:   config[DB_PASSWORD],
		source:     config[DB_SOURCE],
		database:   config[DB_DATABASE],
		collection: config[DB_COLLECTION],
		timeout:    timeout,
	}

	err = instance.Init()

	if err != nil {
		log.Println(err)
		return nil
	}

	return &instance
}

func (db *MongoDB) Init() error {
	var err error

	db.session, err = db.Dial()
	if err != nil {
		return err
	}

	db.session.SetSafe(&mgo.Safe{})
	db.session.SetMode(mgo.Monotonic, true)
	db.session.SetSocketTimeout(1 * time.Hour)

	return nil
}

func (db *MongoDB) Dial() (*mgo.Session, error) {
	if db.url != "" {
		return mgo.Dial(db.url)
	}

	return mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    db.addrs,
		Username: db.username,
		Password: db.password,
		Source:   db.source,
		Timeout:  time.Duration(db.timeout) * time.Millisecond,
	})
}

func (db *MongoDB) Recover() {
	go func() {
		if isLocked {
			return
		}

		isLocked = true
		session, _ := db.Dial()
		if db.session != nil {
			db.session.Close()
		}

		db.session = session
		isLocked = false
	}()
}

func (db *MongoDB) WaitRecover() {
	db.session.Refresh()
}

func (db *MongoDB) IsAlive() (bool, error) {
	if db.session == nil {
		return false, errors.New("Session is null")
	}

	return db.session.Ping() == nil, nil
}

// Must close session when done
func (db *MongoDB) GetCopySession() *mgo.Session {
	return db.session.Copy()
}

func (db *MongoDB) GetCopyCollection() (collection *mgo.Collection, close func()) {
	session := db.GetCopySession()
	close = session.Close
	collection = session.DB(db.database).C(db.collection)
	return
}

func (db *MongoDB) Database() string {
	return db.database
}

func (db *MongoDB) Collection() string {
	return db.collection
}

func (db *MongoDB) GetSession() *mgo.Session {
	return db.session
}

func (db *MongoDB) GetCollection() *mgo.Collection {
	return db.GetDatabase().C(db.collection)
}

func (db *MongoDB) GetDatabase() *mgo.Database {
	return db.GetSession().DB(db.database)
}

func (db *MongoDB) CreateIndexIfNotExist(collection string, index mgo.Index) {
	dbSession := db.session.Copy()
	defer dbSession.Close()

	col := dbSession.DB(db.database).C(collection)
	indexes, e := col.Indexes()
	if e != nil {
		return
	}

	for _, item := range indexes {
		if item.Name == index.Name {
			return
		}
	}

	e = col.EnsureIndex(index)
}

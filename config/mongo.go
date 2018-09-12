package config

import (
	"github.com/BurntSushi/toml"
	"github.com/globalsign/mgo"
	"log"
)

type Mongo struct {
	Server   string
	Port     string
	Database string
	Username string
	Password string
}

func (c *Mongo) Read() {
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Fatal(err)
	}
}

func Connect() *mgo.Database {
	mongoConfig := Mongo{}
	mongoConfig.Read()

	mongoDBDialinfo := &mgo.DialInfo{
		Addrs:    []string{mongoConfig.Server + ":" + mongoConfig.Port},
		Database: mongoConfig.Database,
		Username: mongoConfig.Username,
		Password: mongoConfig.Password,
	}

	session, err := mgo.DialWithInfo(mongoDBDialinfo)
	if err != nil {
		log.Fatal(err)
	}

	return session.DB(mongoConfig.Database)
}

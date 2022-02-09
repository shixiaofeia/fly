package mongo

import (
	"log"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name  string
	Phone string
}

func TestInit(t *testing.T) {
	var err error

	if err = Init(Conf{Address: "127.0.0.1:27017"}); err != nil {
		log.Fatal(err)
	}

	c := NewCollection("test", "people")
	if err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"}); err != nil {
		log.Fatal(err)
	}

	if err = c.Update(bson.M{"name": "Ale"}, bson.M{"$set": bson.M{"phone": "123456"}}); err != nil {
		log.Fatal(err)
	}

	result := Person{}
	if err = c.Find(bson.M{"name": "Ale"}).One(&result); err != nil {
		log.Fatal(err)
	}
	log.Println("Phone:", result.Phone)

	if _, err = c.RemoveAll(bson.M{"name": "Ale"}); err != nil {
		log.Fatal(err)
	}

	countNum, err := c.Count()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(countNum)
}

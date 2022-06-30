package mongo

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Name  string
	Phone string
}

func TestNewCollection(t *testing.T) {
	var (
		err  error
		conf = Conf{
			Host:        "127.0.0.1",
			Port:        "27017",
			DataBase:    "file",
			MaxPoolSize: 10,
		}
		recordM = &Person{ID: bson.NewObjectId(), Name: "fly", Phone: "+86 123456"}
	)

	if err = Init(conf); err != nil {
		t.Fatalf("init err: %v", err)
	}

	cli := NewCollection("people")
	if err = cli.Insert(recordM); err != nil {
		t.Fatalf("insert err: %v", err)
	}

	if err = cli.Update(bson.M{"name": "fly"}, bson.M{"$set": bson.M{"phone": "+86 23456789"}}); err != nil {
		t.Fatalf("update err: %v", err)
	}

	result := Person{}
	if err = cli.Find(bson.M{"_id": recordM.ID}).One(&result); err != nil {
		t.Fatalf("find err: %v", err)
	}

	t.Logf("update after phone: %s", result.Phone)

	if _, err = cli.RemoveAll(bson.M{"name": "fly"}); err != nil {
		t.Fatalf("remove err: %v", err)
	}

	countNum, err := cli.Count()
	if err != nil {
		t.Fatalf("count err: %v", err)
	}
	t.Logf("count num: %d", countNum)
}

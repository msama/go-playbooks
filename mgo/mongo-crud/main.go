package main

import (
	"log"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const _MONGO_TEST_URL = "mongodb://localhost/playbooks_crud"
const _DATABASE = "playbooks_crud"
const _COLLECTION = "documents"

// A nested struct used to obser how to deal with
// complex types
type Nested struct {
	Name  string `bson:"name",json:"name"`
	Count int    `bson:"count",json:"count"`
}

// The document on which we will perform our CRUD operations
type Document struct {
	Id     bson.ObjectId `bson:"_id",json:"id"`
	Nested *Nested       `bson:"nested",json:"nested"`
}

func Create(session *mgo.Session, doc *Document) error {
	collection := session.DB(_DATABASE).C(_COLLECTION)
	return collection.Insert(doc)
}

func Read(session *mgo.Session, id bson.ObjectId) (*Document, error) {
	collection := session.DB(_DATABASE).C(_COLLECTION)

	doc := &Document{}
	if err := collection.Find(bson.M{"_id": id}).One(&doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func Update(session *mgo.Session, doc *Document) error {
	collection := session.DB(_DATABASE).C(_COLLECTION)
	return collection.Update(bson.M{"_id": doc.Id}, doc)
}

func Delete(session *mgo.Session, doc *Document) error {
	collection := session.DB(_DATABASE).C(_COLLECTION)
	return collection.Remove(bson.M{"_id": doc.Id})
}

func main() {
	masterSession, err := mgo.Dial(_MONGO_TEST_URL)
	if err != nil {
		log.Fatalf("Could not connect to %s\n", _MONGO_TEST_URL)
	}
	defer masterSession.Close()
	masterSession.SetSafe(&mgo.Safe{WMode: "majority"})

	doc := &Document{
		Id: bson.NewObjectId(),
		Nested: &Nested{
			Name:  "foo",
			Count: 1,
		},
	}

	if err := Create(masterSession, doc); err != nil {
		log.Fatalf("Error crating document %v\n", err)
	}
	log.Printf("Created %+v\n", doc.Nested)

	load, err := Read(masterSession, doc.Id)
	if err != nil {
		log.Fatalf("Error reading document %v\n", err)
	}
	log.Printf("Read %+v\n", load.Nested)

	doc.Nested.Name = "bar"
	doc.Nested.Count = 9

	if err := Update(masterSession, doc); err != nil {
		log.Fatalf("Error updating document %v\n", err)
	}
	log.Printf("Updated %+v\n", doc.Nested)

	load, err = Read(masterSession, doc.Id)
	if err != nil {
		log.Fatalf("Error reading document %v\n", err)
	}
	log.Printf("Read %+v\n", load.Nested)

	if err := Delete(masterSession, doc); err != nil {
		log.Fatalf("Error deleting document %v\n", err)
	}
	log.Printf("Deleted\n")
}

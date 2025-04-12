package types

import "go.mongodb.org/mongo-driver/bson"

type BsonM bson.M

func Empty() BsonM {
	return make(BsonM)
}

func (b BsonM) SetID(id string) BsonM {
	b["_id"] = id
	return b
}

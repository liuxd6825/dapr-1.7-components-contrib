package model

import "github.com/google/uuid"

type ObjectID string

var NilObjectID = ObjectID("")

func ObjectIDFromHex(str string) (ObjectID, error) {
	return ObjectID(str), nil
}

func NewObjectID() ObjectID {
	return ObjectID(uuid.New().String())
}

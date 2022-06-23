package model

import "github.com/google/uuid"

type ObjectID string

var NilObjectID = ""

func ObjectIDFromHex(str string) (string, error) {
	return str, nil
}

func NewObjectID() string {
	return uuid.New().String()
}

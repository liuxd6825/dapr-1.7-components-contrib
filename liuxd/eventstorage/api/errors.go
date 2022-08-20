package api

import (
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

func NotDuplicateKeyError(err error) (error, bool) {
	if ok := IsDuplicateKeyError(err); ok {
		return nil, ok
	}
	return err, false
}

func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	} else {
		msg := err.Error()
		if strings.Contains(msg, "E11000 duplicate") && strings.Contains(msg, "_id:") {
			return true
		}
	}
	return mongo.IsDuplicateKeyError(err)
}

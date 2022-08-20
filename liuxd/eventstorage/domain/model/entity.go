package model

import "github.com/google/uuid"

type Entity interface {
	GetId() string
	GetTenantId() string
}

func NewObjectID() string {
	return uuid.New().String()
}

package model

import "github.com/google/uuid"

type Entity interface {
	GetId() string
	SetId(v string)
	GetTenantId() string
}

func NewObjectID() string {
	return uuid.New().String()
}

type EntityBuilder interface {
	NewEntity() interface{}
	NewEntities() interface{}
}

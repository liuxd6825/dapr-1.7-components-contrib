package service

import (
	ctx "context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
	"testing"
)

func Test_RelationServiceCreate(t *testing.T) {
	mongodb, err := newTestMongoDb()
	if err != nil {
		t.Error(err)
		return
	}
	service := NewRelationService(mongodb)
	relation := &model.RelationEntity{
		Id:          model.NewObjectID(),
		TableName:   "test_relation",
		AggregateId: model.NewObjectID(),
		Items:       map[string]string{},
	}
	relation.AddItem("CaseId", "caseId")
	relation.AddItem("TaskId", "taskId")
	relation.AddItem("UserId", "userId")
	err = service.Create(ctx.Background(), relation)
	if err != nil {
		t.Error(err)
	}
}

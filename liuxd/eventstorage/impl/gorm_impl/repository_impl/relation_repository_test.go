package repository_impl

import (
	ctx "context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"testing"
)

func Test_RelationServiceCreate(t *testing.T) {
	gormDb, err := newGormDB()
	if err != nil {
		t.Error(err)
		return
	}
	service := NewRelationRepository(gormDb)
	relation := &model.Relation{
		Id:          model.NewObjectID(),
		TenantId:    TEST_TENANT_ID,
		TableName:   "test_relation",
		AggregateId: model.NewObjectID(),
		Items:       map[string]string{},
	}
	relation.AddItem("CaseId", "caseId")
	relation.AddItem("TaskId", "taskId")
	relation.AddItem("UserId", "userId")
	err = service.Create(ctx.Background(), TEST_TENANT_ID, relation)
	if err != nil {
		t.Error(err)
	}
}

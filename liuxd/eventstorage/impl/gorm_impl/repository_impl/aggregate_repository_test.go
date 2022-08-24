package repository_impl

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/impl/gorm_impl/db"
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

const (
	TEST_TENANT_ID = "test"
)

func TestAggregateRepository_Create(t *testing.T) {
	gormDb, err := newGormDB()
	if err != nil {
		t.Error(err)
		return
	}

	repos := NewAggregateRepository(gormDb)
	agg := &model.Aggregate{
		Id:             newId(),
		TenantId:       TEST_TENANT_ID,
		AggregateId:    newId(),
		AggregateType:  "type",
		SequenceNumber: 1,
	}
	ctx := context.Background()
	err = db.NewSession(gormDb).UseTransaction(ctx, func(ctx context.Context) error {
		err := repos.Create(ctx, agg)
		if err != nil {
			return err
		}

		t.Logf("AggreageId = %v", agg.Id)

		if fAgg, ok, err := repos.FindById(ctx, TEST_TENANT_ID, agg.Id); err != nil {
			return err
		} else if ok {
			println(fmt.Sprintf("FindById() agg = %v", fAgg))
		}

		agg.SequenceNumber = agg.SequenceNumber + 1
		if err := repos.Update(ctx, agg); err != nil {
			return err
		}

		if agg, ok, sn, err := repos.NextSequenceNumber(ctx, TEST_TENANT_ID, agg.AggregateId, 1); err != nil {
			return err
		} else if ok {
			println(fmt.Sprintf("NextSequenceNumber() agg =%v", agg))
			println(fmt.Sprintf("NextSequenceNumber() sn=%v", sn))
		}

		if agg, ok, err := repos.UpdateIsDelete(ctx, TEST_TENANT_ID, agg.Id); err != nil {
			return err
		} else if ok {
			println(fmt.Sprintf("UpdateIsDelete() sn=%v", agg))
		}

		if err := repos.DeleteById(ctx, TEST_TENANT_ID, agg.Id); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		t.Error(err)
	}

}

func TestAggregateRepository_FindById(t *testing.T) {
	gormDb, err := newGormDB()
	if err != nil {
		t.Error(err)
		return
	}
	repos := NewAggregateRepository(gormDb)
	ctx := context.Background()
	fAgg, ok, err := repos.FindById(ctx, TEST_TENANT_ID, "485ef5fb-5b97-4109-93f2-713f28aaa782")
	if err != nil {
		t.Error(err)
		return
	} else if ok {
		println(fAgg)
	}
	return
}

func newId() string {
	return uuid.New().String()
}

func newGormDB() (*gorm.DB, error) {
	gormDb, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:11111111@tcp(127.0.0.1:3306)/dapr_es?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                               // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                              // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                              // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                              // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                             // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	return gormDb, err
}

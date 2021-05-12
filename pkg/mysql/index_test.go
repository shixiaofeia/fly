package mysql

import (
	"context"
	"testing"
	"time"
)

func TestGetDB(t *testing.T) {
	if err := Init(Conf{}, Conf{}); err != nil {
		return
	}
	type Fly struct {
		Id         uint32 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
		Name       string `gorm:"column:name;default:'';comment:'姓名'" json:"name"`
		CreateTime int64  `gorm:"column:create_time;default:0" json:"createTime"`
		UpdateTime int64  `gorm:"column:update_time;default:0" json:"updateTime"`
	}
	tx := NewWriteDB(context.TODO())
	if err := tx.AutoMigrate(Fly{}); err != nil {
		t.Error(err)
	}
	flyM := &Fly{Name: "fly"}
	if err := tx.Create(flyM).Error; err != nil {
		t.Error(err)
	}
	if err := tx.Take(flyM).Error; err != nil {
		t.Error(err)
	}
	time.Sleep(2 * time.Second)
	flyM.Name = "fly2"
	if err := tx.Updates(flyM).Error; err != nil {
		t.Error(err)
	}
	// 批量插入钩子函数只对第一条数据生效
	if err := tx.Create(&[]Fly{{Name: "1"}, {Name: "2"}}).Error; err != nil {
		t.Error(err)
	}
}

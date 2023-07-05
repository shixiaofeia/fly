package service

import (
	"fly/example/api/example/model"
	"fly/internal/domain/sqldb"
	"fly/internal/httpcode"
	"fly/pkg/libs/structf"
	"fmt"
	"gorm.io/gorm"
)

type DemoService struct {
	orm *gorm.DB
}

func NewDemoService(orm *gorm.DB) *DemoService {
	return &DemoService{orm: orm}
}

// DemoCreate demo创建.
func (slf *DemoService) DemoCreate(req *model.DemoCreateReq) error {
	var (
		recordM = new(sqldb.Demo)
	)
	structf.Assign(req, recordM)

	if _, err := sqldb.NewDemoSearch(slf.orm).Create(recordM); err != nil {
		return fmt.Errorf("demo create err: %v", err)
	}

	return nil
}

// DemoRecords demo列表.
func (slf *DemoService) DemoRecords(req *model.DemoRecordReq) (res model.DemoRecordResp, err error) {
	res.PageNum, res.PageSize = req.PageNum, req.PageSize
	res.List = make([]*model.DemoRecordItem, 0)

	records, total, err := sqldb.NewDemoSearch(slf.orm).SetName(req.Name).SetPager(req.PageNum, req.PageSize).Find()
	if err != nil {
		return res, fmt.Errorf("demo records err: %v", err)
	}

	for _, v := range records {
		recordM := new(model.DemoRecordItem)
		structf.Assign(v, recordM)
		res.List = append(res.List, recordM)
	}

	res.Total = total

	return res, nil
}

// DemoInfo 单条查询.
func (slf *DemoService) DemoInfo(req *model.DemoInfoReq) (res model.DemoInfoResp, err error) {
	recordM, err := sqldb.NewDemoSearch(slf.orm).SetID(req.DemoID).First()
	if err != nil {
		return res, fmt.Errorf("demo info err: %v", err)
	}

	res.Demo = *recordM

	return res, nil
}

// DemoUpdate 更新.
func (slf *DemoService) DemoUpdate(req *model.DemoUpdateReq) (httpcode.ErrCode, error) {
	var upMap = make(map[string]interface{})

	if req.Name != "" {
		upMap["name"] = req.Name
	}
	if req.IsFree > 0 {
		upMap["is_free"] = req.IsFree
	}
	if req.Amount.IsPositive() {
		upMap["amount"] = req.Amount
	}
	if req.Remark != "" {
		upMap["remark"] = req.Remark
	}
	if len(upMap) == 0 {
		return httpcode.ParamErr, fmt.Errorf("upmap is nil")
	}

	if err := sqldb.NewDemoSearch(slf.orm).SetID(req.DemoID).UpdateByMap(upMap); err != nil {
		return httpcode.ServiceErr, fmt.Errorf("demo update err: %v", err)
	}

	return httpcode.Code200, nil
}

// DemoDelete 删除.
func (slf *DemoService) DemoDelete(req *model.DemoDeleteReq) error {
	if err := sqldb.NewDemoSearch(slf.orm).SetID(req.DemoID).Delete(); err != nil {
		return fmt.Errorf("demo delete err: %v", err)
	}

	return nil
}

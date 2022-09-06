package service

import (
	"fly/example/api/example/model"
	"fly/internal/domain/sqldb"
	"fly/pkg/httpcode"
	"fly/pkg/libs/structf"
	"fmt"
)

type DemoService struct{}

func NewDemoService() *DemoService {
	return &DemoService{}
}

// DemoCreate demo创建.
func (slf *DemoService) DemoCreate(req *model.DemoCreateReq) error {
	var (
		factory = sqldb.NewDemoSearch(nil)
		recordM = new(sqldb.Demo)
	)
	structf.Assign(req, recordM)

	if _, err := factory.Create(recordM); err != nil {
		return fmt.Errorf("demo create err: %v", err)
	}

	return nil
}

// DemoRecords demo列表.
func (slf *DemoService) DemoRecords(req *model.DemoRecordReq, res *model.DemoRecordResp) error {
	var (
		factory = sqldb.NewDemoSearch(nil)
	)
	structf.Assign(req, res)
	structf.Assign(req, factory)

	records, total, err := factory.Find()
	if err != nil {
		return fmt.Errorf("demo records err: %v", err)
	}

	for _, v := range records {
		recordM := new(model.DemoRecordItem)
		structf.Assign(v, recordM)
		res.List = append(res.List, recordM)
	}

	res.Total = total

	return nil
}

// DemoInfo 单条查询.
func (slf *DemoService) DemoInfo(req *model.DemoInfoReq, res *model.DemoInfoResp) error {
	var (
		factory = sqldb.NewDemoSearch(nil)
	)
	factory.Id = req.DemoId

	recordM, err := factory.First()
	if err != nil {
		return fmt.Errorf("demo info err: %v", err)
	}

	structf.Assign(recordM, res)

	return nil
}

// DemoUpdate 更新.
func (slf *DemoService) DemoUpdate(req *model.DemoUpdateReq) (httpcode.ErrCode, error) {
	var (
		factory = sqldb.NewDemoSearch(nil)
		upMap   = make(map[string]interface{})
	)
	factory.Id = req.DemoId

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

	if err := factory.UpdateByMap(upMap); err != nil {
		return httpcode.ServiceErr, fmt.Errorf("demo update err: %v", err)
	}

	return httpcode.Code200, nil
}

// DemoDelete 删除.
func (slf *DemoService) DemoDelete(req *model.DemoDeleteReq) error {
	var (
		factory = sqldb.NewDemoSearch(nil)
	)

	factory.Id = req.DemoId

	if err := factory.Delete(); err != nil {
		return fmt.Errorf("demo delete err: %v", err)
	}

	return nil
}

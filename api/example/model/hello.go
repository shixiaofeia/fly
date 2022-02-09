package model

type (
	HelloReq struct {
		Name string `json:"name" validate:"nonzero"`
	}
	HelloRes struct {
		Name string `json:"name"`
	}

	ExportRes struct {
		Id   int    `json:"id"`
		Name string `json:"name" excel:"-"`
		Age  int    `json:"age"`
	}
)

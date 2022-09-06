package types

type (
	Pager struct {
		PageNum  int   `json:"pageNum"`
		PageSize int   `json:"pageSize"`
		Total    int64 `json:"total"`
	}
)

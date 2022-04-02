package constvar

type (
	Pager struct {
		Page  int   `json:"page"`
		Size  int   `json:"size"`
		Total int64 `json:"total"`
	}
)

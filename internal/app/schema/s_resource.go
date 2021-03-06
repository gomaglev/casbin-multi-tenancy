package schema

// Resource 租户对象
type Resource struct {
	Lang string `json:"lang"` // 语言
	ID   string `json:"id"`   // 名称
	Type string `json:"type"` // Type
}

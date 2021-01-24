package model

import "github.com/aaronchen2k/tester/internal/server/model/base"

type Iso struct {
	BaseModel
	base.TestEnv

	Name string
	Path string
	Size int

	ResolutionHeight  int
	ResolutionWidth   int
	suggestDiskSize   int
	suggestMemorySize int
}

func (Iso) TableName() string {
	return "biz_iso"
}

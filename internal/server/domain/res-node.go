package domain

import _const "github.com/aaronchen2k/tester/internal/pkg/const"

type ResNode struct {
	Id       string         `json:"id"`
	Name     string         `json:"name"`
	Type     _const.ResType `json:"type"`
	Key      string         `json:"key"`
	Children []*ResNode     `json:"children"`

	HostId string `json:"hostId,omitempty"`
	NodeId string `json:"nodeId,omitempty"`
}
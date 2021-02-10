package model

import (
	_const "github.com/aaronchen2k/tester/internal/pkg/const"
)

type Node struct {
	BaseModel

	Cluster   string `json:"cluster"`
	InstCount int    `gorm:"default:0" json:"instCount"`
	TaskCount int    `gorm:"-" json:"taskCount"`

	Ident string `json:"ident"`
	Type  string `json:"type"`
	Name  string `json:"name"`

	Ip      string `json:"ip"`
	Port    int    `son:"port"`
	WorkDir string `json:"workDir,omitempty"`

	Username string `json:"username"`
	Password string `json:"password"`

	SshPort int               `json:"sshPort,omitempty"`
	VncPort int               `json:"vncPort,omitempty"`
	Status  _const.HostStatus `json:"status"`
}

func NewNode() Node {
	node := Node{}
	return node
}

func (Node) TableName() string {
	return "biz_node"
}

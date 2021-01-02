package _domain

import (
	_const "github.com/aaronchen2k/openstc/internal/pkg/libs/const"
	"time"
)

type Vm struct {
	Id                int
	Name              string
	DiskSize          int
	MemorySize        int
	CdromSys          string
	CdromDriver       string
	DefPath           string
	ImagePath         string
	BackingImagePath  string
	WorkDir           string
	PublicIp          string
	PublicPort        int
	MacAddress        string
	ResolutionHeight  int
	ResolutionWidth   int
	RpcPort           int
	SshPort           int
	VncPort           int
	DestroyAt         time.Time
	FirstDetectedTime time.Time
	HostId            int
	BackingImageId    int
	Status            _const.VmStatus
}

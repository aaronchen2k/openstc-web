package _domain

import (
	_const "github.com/aaronchen2k/tester/internal/pkg/const"
	"time"
)

type BuildTo struct {
	WorkDir    string
	ProjectDir string
	AppPath    string

	ID           uint
	Serial       string
	Priority     int
	ComputerIp   string
	ComputerPort int
	DeviceIp     string

	BuildType   _const.BuildType
	BrowserType _const.BrowserType
	BrowserVer  string
	AppiumPort  int

	QueueId uint

	ScriptUrl   string
	ScmAddress  string
	ScmAccount  string
	ScmPassword string

	AppUrl          string
	BuildCommands   string
	ResultFiles     string
	KeepResultFiles MyBool
	ResultPath      string
	ResultMsg       string

	StartTime    time.Time
	CompleteTime time.Time

	Progress _const.BuildProgress
	Status   _const.BuildStatus
}

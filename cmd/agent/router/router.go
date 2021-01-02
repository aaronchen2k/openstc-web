package router

import (
	"fmt"
	"github.com/aaronchen2k/openstc/cmd/agent/router/handler"
	"github.com/aaronchen2k/openstc/internal/agent/cfg"
	_const "github.com/aaronchen2k/openstc/internal/pkg/libs/const"
	_logUtils "github.com/aaronchen2k/openstc/internal/pkg/libs/log"
	"github.com/smallnest/rpcx/server"
	"strconv"
)

func App() {
	addr := agentConf.Inst.Ip + ":" + strconv.Itoa(agentConf.Inst.Port)

	srv := server.NewServer()

	if agentConf.Inst.Platform == _const.Android || agentConf.Inst.Platform == _const.Ios {
		srv.RegisterName("appium", new(handler.AppiumAction), "")
	} else if agentConf.Inst.Platform == _const.Host {
		srv.RegisterName("vm", new(handler.VmAction), "")
		srv.RegisterName("image", new(handler.ImageAction), "")
	} else if agentConf.Inst.Platform == _const.Vm {
		srv.RegisterName("selenium", new(handler.SeleniumAction), "")
	}

	_logUtils.Info(fmt.Sprintf("start server on %s ...", addr))
	err := srv.Serve("tcp", addr)
	if err != nil {
		_logUtils.Infof("fail to start server on %s, err is %s", err.Error())
	}
}

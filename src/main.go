//go:generate go-bindata -prefix "./www/dist" -fs  ./www/dist/...
package main

import (
	"flag"
	"fmt"
	"github.com/aaronchen2k/openstc/src/libs/common"
	"github.com/aaronchen2k/openstc/src/server"
	"os"
)

var (
	version      = "master"
	configPath   = flag.String("c", "", "配置路径")
	printVersion = flag.Bool("v", false, "打印版本号")
	seederData   = flag.Bool("s", false, "填充基础数据")
	syncPerms    = flag.Bool("p", true, "同步权限")
	printRouter  = flag.Bool("r", false, "打印路由列表")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] [command]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  -c <path>\n")
		fmt.Fprintf(os.Stderr, "    设置配置文件路径\n")
		fmt.Fprintf(os.Stderr, "  -v <true or false> 默认为: false\n")
		fmt.Fprintf(os.Stderr, "    打印版本号\n")
		fmt.Fprintf(os.Stderr, "  -s <true or false> 默认为: false\n")
		fmt.Fprintf(os.Stderr, "    填充基础数据\n")
		fmt.Fprintf(os.Stderr, "  -p <true or false> 默认为: true\n")
		fmt.Fprintf(os.Stderr, "    同步权限\n")
		fmt.Fprintf(os.Stderr, "  -r <true or false> 默认为: false\n")
		fmt.Fprintf(os.Stderr, "    打印路由列表\n")
		fmt.Fprintf(os.Stderr, "\n")
		//flag.PrintDefaults()
	}
	flag.Parse()

	common.InitConfig(*configPath)
	server.Init(version, printVersion, seederData, syncPerms, printRouter)
}

package server

import (
	"fmt"
	"github.com/aaronchen2k/tester/cmd/server/router"
	"github.com/aaronchen2k/tester/cmd/server/router/handler"
	"github.com/aaronchen2k/tester/internal/pkg/utils"
	"github.com/aaronchen2k/tester/internal/server/biz/middleware"
	middlewareUtils "github.com/aaronchen2k/tester/internal/server/biz/middleware/misc"
	"github.com/aaronchen2k/tester/internal/server/biz/redis"
	"github.com/aaronchen2k/tester/internal/server/cfg"
	serverCron "github.com/aaronchen2k/tester/internal/server/cron"
	"github.com/aaronchen2k/tester/internal/server/db"
	"github.com/aaronchen2k/tester/internal/server/model"
	"github.com/aaronchen2k/tester/internal/server/repo"
	"github.com/aaronchen2k/tester/internal/server/service"
	"github.com/facebookgo/inject"
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"

	"github.com/kataras/iris/v12/context"
)

func Init(version string, printVersion, printRouter *bool) {

	db.InitDB()

	// irisServer := server.NewServer(AssetFile()) // 加载前端文件
	irisServer := NewServer(nil)
	if irisServer == nil {
		panic("Http 初始化失败")
	}
	irisServer.App.Logger().SetLevel(serverConf.Config.LogLevel)

	//if common.Config.BinData {
	//	irisServer.App.RegisterView(iris.HTML(irisServer.AssetFile, ".html"))
	//	irisServer.App.HandleDir("/", irisServer.AssetFile)
	//}

	router := router.NewRouter(irisServer.App)
	injectObj(router)
	router.InitService.Init()
	router.App()

	if serverConf.Config.Redis.Enable {
		redisUtils.InitRedisCluster(serverConf.GetRedisUris(), serverConf.Config.Redis.Pwd)
	}

	iris.RegisterOnInterrupt(func() {
		defer db.GetInst().Close()
	})

	// deal with the command
	if *printVersion {
		fmt.Println(fmt.Sprintf("版本号：%s", version))
	}

	if *printRouter {
		fmt.Println("系统权限：")
		fmt.Println()
		routes := irisServer.GetRoutes()
		for _, route := range routes {
			fmt.Println("+++++++++++++++")
			fmt.Println(fmt.Sprintf("名称 ：%s ", route.DisplayName))
			fmt.Println(fmt.Sprintf("路由地址 ：%s ", route.Name))
			fmt.Println(fmt.Sprintf("请求方式 ：%s", route.Act))
			fmt.Println()
		}
	}

	if _utils.IsPortInUse(serverConf.Config.Port) {
		panic(fmt.Sprintf("端口 %d 已被使用", serverConf.Config.Port))
	}

	// start the service
	err := irisServer.Serve()
	if err != nil {
		panic(err)
	}
}

func injectObj(router *router.Router) {
	// inject
	var g inject.Graph
	g.Logger = logrus.StandardLogger()

	if err := g.Provide(
		// db
		&inject.Object{Value: db.GetInst().DB()},

		// repo
		&inject.Object{Value: repo.NewEnvRepo()},

		&inject.Object{Value: repo.NewBuildRepo()},
		&inject.Object{Value: repo.NewCommonRepo()},
		&inject.Object{Value: repo.NewDeviceRepo()},
		&inject.Object{Value: repo.NewExecRepo()},

		&inject.Object{Value: repo.NewIsoRepo()},
		&inject.Object{Value: repo.NewPermRepo()},
		&inject.Object{Value: repo.NewQueueRepo()},
		&inject.Object{Value: repo.NewRoleRepo()},
		&inject.Object{Value: repo.NewTaskRepo()},
		&inject.Object{Value: repo.NewTokenRepo()},
		&inject.Object{Value: repo.NewUserRepo()},

		&inject.Object{Value: repo.NewClusterRepo()},
		&inject.Object{Value: repo.NewComputerRepo()},
		&inject.Object{Value: repo.NewVmTemplRepo()},
		&inject.Object{Value: repo.NewVmRepo()},
		&inject.Object{Value: repo.NewContainerImageRepo()},
		&inject.Object{Value: repo.NewContainerRepo()},

		// middleware
		&inject.Object{Value: middlewareUtils.NewEnforcer()},
		&inject.Object{Value: middleware.NewJwtService()},
		&inject.Object{Value: middleware.NewTokenService()},
		&inject.Object{Value: middleware.NewCasbinService()},

		// service
		&inject.Object{Value: serverCron.NewServerCron()},
		&inject.Object{Value: service.NewEnvService()},

		&inject.Object{Value: service.NewAppiumService()},
		&inject.Object{Value: service.NewBuildService()},
		&inject.Object{Value: service.NewCommonService()},
		&inject.Object{Value: service.NewDeviceService()},
		&inject.Object{Value: service.NewExecService()},
		&inject.Object{Value: service.NewClusterService()},

		&inject.Object{Value: service.NewIsoService()},
		&inject.Object{Value: service.NewPermService()},
		&inject.Object{Value: service.NewQueueService()},
		&inject.Object{Value: service.NewRoleService()},
		&inject.Object{Value: service.NewSeeder()},
		&inject.Object{Value: service.NewSeleniumService()},
		&inject.Object{Value: service.NewTaskService()},
		&inject.Object{Value: service.NewUserService()},

		&inject.Object{Value: service.NewResService()},
		&inject.Object{Value: service.NewDockerImageService()},
		&inject.Object{Value: service.NewContainerService()},
		&inject.Object{Value: service.NewVmTemplService()},
		&inject.Object{Value: service.NewVmService()},

		&inject.Object{Value: service.NewRpcService()},

		// controller
		&inject.Object{Value: handler.NewEnvCtrl()},

		&inject.Object{Value: handler.NewAccountCtrl()},
		&inject.Object{Value: handler.NewAppiumCtrl()},
		&inject.Object{Value: handler.NewDeviceCtrl()},
		&inject.Object{Value: handler.NewFileCtrl()},

		&inject.Object{Value: handler.NewInitCtrl()},
		&inject.Object{Value: handler.NewPermCtrl()},
		&inject.Object{Value: handler.NewUserCtrl()},

		&inject.Object{Value: handler.NewTaskCtrl()},

		&inject.Object{Value: handler.NewRoleCtrl()},
		&inject.Object{Value: handler.NewBuildCtrl()},

		&inject.Object{Value: handler.NewClusterCtrl()},
		&inject.Object{Value: handler.NewMachineCtrl()},
		&inject.Object{Value: handler.NewContainerImageCtrl()},
		&inject.Object{Value: handler.NewContainerCtrl()},
		&inject.Object{Value: handler.NewVmTemplCtrl()},
		&inject.Object{Value: handler.NewVmCtrl()},

		// router
		&inject.Object{Value: router},
	); err != nil {
		logrus.Fatalf("provide usecase objects to the Graph: %v", err)
	}

	err := g.Populate()
	if err != nil {
		logrus.Fatalf("populate the incomplete Objects: %v", err)
	}
}

type Server struct {
	App       *iris.Application
	AssetFile http.FileSystem
}

func NewServer(assetFile http.FileSystem) *Server {
	app := iris.Default()
	return &Server{
		App:       app,
		AssetFile: assetFile,
	}
}

func (s *Server) Serve() error {
	if serverConf.Config.Https {
		host := fmt.Sprintf("%s:%d", serverConf.Config.Host, 443)
		if err := s.App.Run(iris.TLS(host, serverConf.Config.CertPath, serverConf.Config.CertKey)); err != nil {
			return err
		}
	} else {
		if err := s.App.Run(
			iris.Addr(fmt.Sprintf("%s:%d", serverConf.Config.Host, serverConf.Config.Port)),
			iris.WithoutServerError(iris.ErrServerClosed),
			iris.WithOptimizations,
			iris.WithTimeFormat(time.RFC3339),
		); err != nil {
			return err
		}
	}

	return nil
}

type PathName struct {
	Name   string
	Path   string
	Method string
}

// 获取路由信息
func (s *Server) GetRoutes() []*model.Permission {
	var rrs []*model.Permission
	names := getPathNames(s.App.GetRoutesReadOnly())
	if serverConf.Config.Debug {
		fmt.Println(fmt.Sprintf("路由权限集合：%v", names))
		fmt.Println(fmt.Sprintf("Iris App ：%v", s.App))
	}
	for _, pathName := range names {
		if !isPermRoute(pathName.Name) {
			rr := &model.Permission{Name: pathName.Path, DisplayName: pathName.Name, Description: pathName.Name, Act: pathName.Method}
			rrs = append(rrs, rr)
		}
	}
	return rrs
}

func getPathNames(routeReadOnly []context.RouteReadOnly) []*PathName {
	var pns []*PathName
	if serverConf.Config.Debug {
		fmt.Println(fmt.Sprintf("routeReadOnly：%v", routeReadOnly))
	}
	for _, s := range routeReadOnly {
		pn := &PathName{
			Name:   s.Name(),
			Path:   s.Path(),
			Method: s.Method(),
		}
		pns = append(pns, pn)
	}

	return pns
}

// 过滤非必要权限
func isPermRoute(name string) bool {
	exceptRouteName := []string{"OPTIONS", "GET", "POST", "HEAD", "PUT", "PATCH", "payload"}
	for _, er := range exceptRouteName {
		if strings.Contains(name, er) {
			return true
		}
	}
	return false
}

package service

import (
	"fmt"
	_const "github.com/aaronchen2k/tester/internal/pkg/const"
	_domain "github.com/aaronchen2k/tester/internal/pkg/domain"
	_stringUtils "github.com/aaronchen2k/tester/internal/pkg/libs/string"
	"github.com/aaronchen2k/tester/internal/server/model"
	"github.com/aaronchen2k/tester/internal/server/repo"
	serverUtils "github.com/aaronchen2k/tester/internal/server/utils/common"
	"math/rand"
	"strings"
)

type VmService struct {
	RpcService *RpcService `inject:""`
	ResService *ResService `inject:""`

	VmRepo       *repo.VmRepo       `inject:""`
	VmTemplRepo  *repo.VmTemplRepo  `inject:""`
	ClusterRepo  *repo.ClusterRepo  `inject:""`
	ComputerRepo *repo.ComputerRepo `inject:""`

	IsoRepo   *repo.IsoRepo   `inject:""`
	QueueRepo *repo.QueueRepo `inject:""`
}

func NewVmService() *VmService {
	return &VmService{}
}

func (s *VmService) CreateByQueue(queue model.Queue) (err error) {
	templ := s.VmTemplRepo.Get(queue.VmTemplId)
	computer := s.ComputerRepo.GetByIndent(templ.Computer)
	cluster := s.ClusterRepo.GetByIdent(templ.Cluster)

	vmName := serverUtils.GenVmHostName(queue.ID, templ.OsPlatform, templ.OsType, templ.OsLang)
	vmIdent, err := s.ResService.CreateVm(vmName, templ, computer, cluster)

	if err != nil || vmIdent == "" { //  fail to create
		return
	}

	vm := model.Vm{
		Name:       vmName,
		Ident:      vmIdent,
		Computer:   computer.Ident,
		Cluster:    computer.Cluster,
		ComputerId: computer.ID,
		ClusterId:  cluster.ID,
		Status:     _const.VmCreated,
	}
	s.VmRepo.Save(&vm) // vm status: created

	queue.VmId = vm.ID
	s.QueueRepo.SetAndLaunchVm(queue)               // queue progress: launch_vm
	s.VmRepo.UpdateStatus(vm.ID, _const.VmLaunched) // vm status: launched

	s.ComputerRepo.AddInstCount(computer.ID)

	return
}

func (s *VmService) Register(vm _domain.Vm) (result _domain.RpcResult) {
	err := s.VmRepo.Register(vm)
	if err != nil {
		result.Fail(fmt.Sprintf("fail to register host %s ", vm.MacAddress))
	}
	return
}

func (s *VmService) genVmName(imageName string) (name string) {
	uuid := strings.Replace(_stringUtils.NewUUID(), "-", "", -1)
	name = strings.Replace(imageName, "backing", uuid, -1)

	return
}

func (s *VmService) genRandomMac() (mac string) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	buf[0] |= 2
	mac = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x\n", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
	return
}

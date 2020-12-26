package repo

import (
	_domain "github.com/aaronchen2k/openstc-common/src/domain"
	_const "github.com/aaronchen2k/openstc-common/src/libs/const"
	"github.com/aaronchen2k/openstc/src/model"
	"gorm.io/gorm"
	"time"
)

func NewVmRepo() *VmRepo {
	return &VmRepo{}
}

type VmRepo struct {
	CommonRepo
	DB *gorm.DB `inject:""`
}

func (r *VmRepo) Register(vm _domain.Vm) (err error) {
	// just update status by mac for exist vm
	r.DB.Model(&vm).Where("mac=?", vm.MacAddress).
		Updates(
			map[string]interface{}{"status": vm.Status, "ip": vm.PublicIp, "port": vm.PublicPort, "workDir": vm.WorkDir,
				"lastRegisterDate": time.Now(), "updatedAt": time.Now()})

	return
}

func (r *VmRepo) GetById(id uint) (vm model.Vm) {
	r.DB.Where("ID=?", id).First(&vm)
	return
}
func (r *VmRepo) GetByMac(mac string) (vm model.Vm) {
	r.DB.Where("mac=?", mac).First(&vm)
	return
}

func (r *VmRepo) Save(po model.Vm) {
	r.DB.Model(&po).Omit("").Create(&po)
	return
}

func (r *VmRepo) Launch(vm _domain.Vm) {
	r.DB.Model(&vm).Where("id=?", vm.Id).
		Updates(
			map[string]interface{}{"status": "launch", "imagePath": vm.ImagePath,
				"defPath": vm.DefPath, "updatedAt": time.Now()})

	return
}

func (r *VmRepo) UpdateStatusByNames(vms []string, status _const.VmStatus) {
	db := r.DB.Model(&model.Vm{}).Where("name = IN (?)", vms)

	if status == "running" {
		db.Where("AND status != 'active'")
	}

	db.Updates(map[string]interface{}{"status": status, "updatedAt": time.Now()})
}

func (r *VmRepo) DestroyMissedVmsStatus(vms []string, hostId uint) {
	db := r.DB.Model(&model.Vm{}).Where("hostId=? AND status!=?", hostId, "destroy")

	if len(vms) > 0 {
		db.Where("AND name NOT IN (?)", vms)
	}

	db.Updates(map[string]interface{}{"status": "destroy", "updatedAt": time.Now()})
}

func (r *VmRepo) FailToCreate(id uint, msg string) {
	r.DB.Model(&model.Vm{}).
		Where("id=?", id).
		Updates(map[string]interface{}{"msg": _const.VmFailToCreate, "updatedAt": time.Now()})
}

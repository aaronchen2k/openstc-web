package repo

import (
	_domain "github.com/aaronchen2k/tester/internal/pkg/domain"
	_commonUtils "github.com/aaronchen2k/tester/internal/pkg/libs/common"
	_const "github.com/aaronchen2k/tester/internal/pkg/libs/const"
	"github.com/aaronchen2k/tester/internal/server/model"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type HostRepo struct {
	CommonRepo
	DB *gorm.DB `inject:""`
}

func NewHostRepo() *HostRepo {
	return &HostRepo{}
}

func (r *HostRepo) Register(host _domain.Host) (po model.Host, err error) {

	err = r.DB.Where("ip = ?", host.Ip).First(&po).Error
	if err != gorm.ErrRecordNotFound {
		err = r.DB.Model(&host).Omit("").Create(&host).Error
		return
	} else {
		err = r.DB.Model(&host).Updates(host).Error
		return
	}
}

func (r *HostRepo) Query(keywords string, pageNo, pageSize int) (hosts []model.Host, total int64, err error) {
	query := r.DB.Select("*")
	if keywords != "" {
		query = query.Where("title LIKE ?", "%"+keywords+"%")
	}
	if pageNo > 0 {
		query = query.Offset((pageNo - 1) * pageSize).Limit(pageSize)
	}

	err = query.Find(&hosts).Error
	err = r.DB.Model(&model.Host{}).Count(&total).Error

	return
}

func (r *HostRepo) Get(id int) (host model.Host) {
	r.DB.Where("id=?", id).First(&host)
	return
}

func (r *HostRepo) QueryByImages(images []int, ids []int) (hostId, backingImageId int) {
	sql := `SELECT r.hostId, r.backingImageId imageId 
			FROM BizHostCapability_relation r 
		    INNER JOIN BizHost host on r.hostId = host.id 

	        WHERE host.status = 'active' 
			AND r.backingImageId IN (` +
		strings.Join(_commonUtils.IntToStrArr(images), ",") +
		`) AND host.id IN ("` +
		strings.Join(_commonUtils.IntToStrArr(ids), ",") +
		`) LIMIT 1`

	r.DB.Raw(sql).Scan(&ids)
	return
}

func (r *HostRepo) QueryIdle(host int) (ret []map[string]int) {
	sql := `SELECT temp.hostId, temp.vmCount 
			FROM (
				SELECT DISTINCT host.id hostId, IFNULL(sub.num, 0) vmCount
				FROM BizHost host
				LEFT JOIN
					(SELECT hostId, COUNT(1) num
						FROM BizVm
						WHERE status != 'destroy' AND NOT deleted AND NOT disabled
						GROUP BY hostId) sub
					ON host.id = sub.hostId
			) temp
			WHERE temp.vmCount <= ` + strconv.Itoa(_const.MaxVmOnHost)

	r.DB.Raw(sql).Scan(&ret)
	return
}
package deviceService

import (
	androidService "github.com/aaronchen2k/tester/internal/agent/service/device/android"
	iosService "github.com/aaronchen2k/tester/internal/agent/service/device/ios"
	"github.com/aaronchen2k/tester/internal/agent/utils/common"
	_domain "github.com/aaronchen2k/tester/internal/pkg/domain"
	_const "github.com/aaronchen2k/tester/internal/pkg/libs/const"
)

func RefreshDevices() []_domain.DeviceInst {
	if agentUntils.IsAndroidAgent() {
		androidService.Devices = androidService.GetDeviceInsts()
		return androidService.Devices
	} else if agentUntils.IsIosAgent() {
		iosService.Devices = iosService.GetDeviceInsts()
		return iosService.Devices
	}

	return nil
}

func GetDevice(serial string) (_domain.DeviceInst, bool) {
	var devices []_domain.DeviceInst
	if agentUntils.IsAndroidAgent() {
		devices = androidService.Devices
	} else if agentUntils.IsIosAgent() {
		devices = iosService.Devices
	}

	for _, dev := range devices {
		if dev.Serial == serial {
			return dev, true
		}
	}

	return _domain.DeviceInst{}, false
}

func IsValid(devs []_domain.DeviceInst, serial string) bool {
	for _, dev := range devs {
		if dev.Serial == serial {
			if dev.DeviceStatus == _const.DeviceActive && dev.AppiumStatus == _const.ServiceActive {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func GetDeviceSerials(devices []_domain.DeviceInst) []string {
	ret := make([]string, 0)

	for _, dev := range devices {
		ret = append(ret, dev.Serial)
	}

	return ret
}

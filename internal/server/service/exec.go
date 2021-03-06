package service

import (
	_const "github.com/aaronchen2k/tester/internal/pkg/const"
	"github.com/aaronchen2k/tester/internal/server/model"
	"github.com/aaronchen2k/tester/internal/server/repo"
)

type ExecService struct {
	ResService      *ResService      `inject:""`
	DeviceService   *DeviceService   `inject:""`
	VmService       *VmService       `inject:""`
	AppiumService   *AppiumService   `inject:""`
	SeleniumService *SeleniumService `inject:""`
	TaskService     *TaskService     `inject:""`

	ExecRepo   *repo.ExecRepo   `inject:""`
	QueueRepo  *repo.QueueRepo  `inject:""`
	DeviceRepo *repo.DeviceRepo `inject:""`
	VmRepo     *repo.VmRepo     `inject:""`
	TaskRepo   *repo.TaskRepo   `inject:""`
}

func NewExecService() *ExecService {
	return &ExecService{}
}

func (s *ExecService) CheckAll() {
	s.SetTimeout()
	s.Run()
	s.Retry()

	s.DestroyTimeout()
}

func (s *ExecService) Run() {
	queuesToBuild := s.QueueRepo.QueryForExec()
	for _, queue := range queuesToBuild {
		s.Exec(queue)
	}
}

func (s *ExecService) Exec(queue model.Queue) {
	if queue.BuildType == _const.SeleniumTest {
		s.SeleniumTest(queue)
	} else if queue.BuildType == _const.AppiumTest {
		s.AppiumTest(queue)
	}
}

func (s *ExecService) SeleniumTest(queue model.Queue) {
	originalProgress := queue.Progress
	var newProgress _const.BuildProgress

	if queue.Progress == _const.ProgressCreated {
		// create kvm
		err := s.VmService.CreateByQueue(queue)
		if err == nil { // success to create
			newProgress = _const.ProgressInProgress
		} else {
			newProgress = _const.ProgressPending
		}

	} else if queue.Progress == _const.ProgressLaunchVm {
		vmId := queue.VmId
		vm := s.VmRepo.GetById(vmId)

		if vm.Status == _const.VmActive { // find ready vm, begin to run test
			result := s.SeleniumService.Run(queue)
			if result.IsSuccess() {
				s.QueueRepo.Run(queue)
				newProgress = _const.ProgressInProgress
			} else { // busy, pending
				s.QueueRepo.Pending(queue.ID)
				newProgress = _const.ProgressPending
			}
		}
	}

	if originalProgress != newProgress { // queue's progress changed
		s.TaskRepo.SetProgress(queue.TaskId, newProgress)
	}
}

func (s *ExecService) AppiumTest(queue model.Queue) {
	serial := queue.Serial
	device := s.DeviceRepo.GetBySerial(serial)

	originalProgress := queue.Progress
	var newProgress _const.BuildProgress

	if s.DeviceService.IsDeviceReady(device) {
		rpcResult := s.AppiumService.Run(queue)

		if rpcResult.IsSuccess() {
			s.QueueRepo.Run(queue) // start
			newProgress = _const.ProgressInProgress
		} else {
			s.QueueRepo.Pending(queue.ID) // pending
			newProgress = _const.ProgressPending
		}
	} else {
		s.QueueRepo.Pending(queue.ID) // pending
		newProgress = _const.ProgressPending
	}

	if originalProgress != newProgress { // progress changed
		s.TaskService.SetProgress(queue.TaskId, newProgress)
	}
}

func (s *ExecService) SetTimeout() {
	queueIds := s.QueueRepo.QueryTimeout()
	s.QueueRepo.SetTimeout(queueIds)
}

func (s *ExecService) Retry() {
	queues := s.QueueRepo.QueryForRetry()

	for _, queue := range queues {
		s.Exec(queue)
	}
}

func (s *ExecService) DestroyTimeout() {
	s.ResService.DestroyTimeout()
}

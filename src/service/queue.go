package service

import (
	_const "github.com/aaronchen2k/openstc-common/src/libs/const"
	"github.com/aaronchen2k/openstc/src/model"
	"github.com/aaronchen2k/openstc/src/repo"
	"strings"
)

type QueueService struct {
	TaskService *TaskService `inject:""`

	QueueRepo  *repo.QueueRepo  `inject:""`
	DeviceRepo *repo.DeviceRepo `inject:""`
}

func NewQueueService() *QueueService {
	return &QueueService{}
}

func (s *QueueService) GenerateFromTask(task model.Task) (count int) {
	if task.BuildType == _const.AppiumTest {
		count = s.GenerateAppiumQueuesFromTask(task)
	} else if task.BuildType == _const.SeleniumTest {
		count = s.GenerateSeleniumQueuesFromTask(task)
	}

	return
}

func (s *QueueService) GenerateAppiumQueuesFromTask(task model.Task) (count int) {
	if len(task.Serials) == 0 {
		return
	}

	var groupId uint
	if task.GroupId != 0 {
		groupId = task.GroupId
	} else {
		groupId = task.ID
	}

	serials := strings.Split(task.Serials, ",")
	for _, serial := range serials {
		serial = strings.TrimSpace(serial)
		if serial == "" {
			continue
		}

		device := s.DeviceRepo.GetBySerial(serial)
		if device.ID != 0 {
			queue := model.NewQueueDetail(serial, task.BuildType, groupId, task.ID, task.Priority,
				"", "", "", "", "",
				task.ScriptUrl, task.ScmAddress, task.ScmAccount, task.ScmPassword,
				task.ResultFiles, task.KeepResultFiles, task.TaskName, task.UserName,
				task.AppUrl, task.BuildCommands)

			s.QueueRepo.Save(&queue)
			count++
		}
	}

	s.QueueRepo.DeleteInSameGroup(task.GroupId, serials) // disable same serial queues in same group

	return
}

func (s *QueueService) GenerateSeleniumQueuesFromTask(task model.Task) (count int) {
	// windows,win10,cn_zh,chrome,84;
	environments := strings.TrimSpace(task.Environments)
	envs := strings.Split(environments, ";")

	if len(envs) == 0 {
		return
	}

	var groupId uint
	if task.GroupId != 0 {
		groupId = task.GroupId
	} else {
		groupId = task.ID
	}

	for _, env := range envs {
		sections := strings.Split(strings.TrimSpace(env), ",")
		if len(sections) < 5 {
			continue
		}

		osPlatform := sections[0]
		osType := sections[1]
		osLang := sections[2]
		browserType := sections[3]
		browserVersion := sections[4]

		queue := model.NewQueueDetail("", task.BuildType, groupId, task.ID, task.Priority,
			_const.OsPlatform(osPlatform), _const.OsType(osType),
			_const.OsLang(osLang), _const.BrowserType(browserType), browserVersion,
			task.ScriptUrl, task.ScmAddress, task.ScmAccount, task.ScmPassword,
			task.ResultFiles, task.KeepResultFiles, task.TaskName, task.UserName,
			"", task.BuildCommands)

		s.QueueRepo.Save(&queue)
		count++
	}

	return
}

func (s *QueueService) SetQueueResult(queueId uint, progress _const.BuildProgress, status _const.BuildStatus) {
	queue := s.QueueRepo.GetQueue(queueId)

	s.QueueRepo.SetQueueStatus(queueId, progress, status)
	s.TaskService.CheckCompleted(queue.TaskId)
}
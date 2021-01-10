package service

import (
	"fmt"
	"github.com/aaronchen2k/tester/internal/pkg/const"
	"github.com/aaronchen2k/tester/internal/pkg/libs/log"
	"github.com/aaronchen2k/tester/internal/server/domain"
	"github.com/aaronchen2k/tester/vendors/github.com/joernott/go-proxmox"
	go_portainer "github.com/aaronchen2k/tester/vendors/github.com/leidruid/go-portainer"
	"strconv"
	"strings"
)

type MachineService struct {
	ClusterService *ClusterService `inject:""`
}

func NewMachineService() *MachineService {
	return &MachineService{}
}

func (s *MachineService) ListVm() (rootNode *domain.ResNode) {
	rootNode = &domain.ResNode{Name: "虚拟机", Type: _const.ResRoot, Id: "0"}
	hosts := s.ClusterService.ListByType("pve")

	for _, host := range hosts {
		id := strconv.Itoa(int(host.ID))

		hostNode := &domain.ResNode{Name: host.Name + "(集群)", Type: _const.ResCluster,
			Id: id, Key: string(_const.ResCluster) + "-" + id}
		rootNode.Children = append(rootNode.Children, hostNode)

		var err error
		address := fmt.Sprintf("%s:%d", host.Ip, host.Port)
		go_proxmox.Proxmox, err = go_proxmox.NewProxMox(address, host.Username, host.Password)
		if err != nil {
			_logUtils.Print("fail to connect proxmox, error: " + err.Error())
			break
		}

		nodes, _ := go_proxmox.Proxmox.Nodes()
		for _, node := range nodes {
			id := node.Id

			nodeNode := &domain.ResNode{Name: node.Node + "(节点)", Type: _const.ResNode,
				Id: id, HostId: hostNode.Id, Key: string(_const.ResNode) + "-" + id}
			hostNode.Children = append(hostNode.Children, nodeNode)

			vmFolderNode := &domain.ResNode{Name: "实例", Type: _const.ResFolder,
				Id: id + "-folder-vms", Key: id + "-folder-vms"}
			nodeNode.Children = append(nodeNode.Children, vmFolderNode)

			templFolderNode := &domain.ResNode{Name: "模板", Type: _const.ResFolder,
				Id: id + "-folder-templs", Key: id + "-folder-templs"}
			nodeNode.Children = append(nodeNode.Children, templFolderNode)

			vms, _ := node.Qemu()
			for _, vm := range vms {
				vmId := strconv.FormatFloat(vm.VMId, 'f', 0, 64)
				isTemplate := false
				if vm.Template == 1 {
					isTemplate = true
				}

				vmNode := &domain.ResNode{Name: vm.Name, Type: _const.ResVm, IsTemplate: isTemplate,
					Id: vmId, HostId: hostNode.Id, NodeId: nodeNode.Id, Key: string(_const.ResVm) + "-" + vmId}

				if !isTemplate {
					vmFolderNode.Children = append(vmFolderNode.Children, vmNode)
				} else {
					templFolderNode.Children = append(templFolderNode.Children, vmNode)
				}
			}
		}
	}

	return
}

func (s *MachineService) ListContainers() (rootNode *domain.ResNode) {
	rootNode = &domain.ResNode{Name: "容器", Type: _const.ResRoot, Id: "0"}
	hosts := s.ClusterService.ListByType("portainer")

	for _, host := range hosts {
		id := strconv.Itoa(int(host.ID))

		hostNode := &domain.ResNode{Name: host.Name + "(集群)", Type: _const.ResCluster,
			Id: id, Key: string(_const.ResCluster) + "-" + id}
		rootNode.Children = append(rootNode.Children, hostNode)

		config := go_portainer.Config{
			Host:     host.Ip,
			Port:     host.Port,
			User:     host.Username,
			Password: host.Password,
			Schema:   "http",
			URL:      "/api",
		}
		portainer := go_portainer.NewPortainer(&config)
		err := portainer.Auth()
		if err != nil {
			_logUtils.Print("fail to connect portainer, error: " + err.Error())
			break
		}

		endpoints, _ := portainer.ListEndpoints()
		for _, endpoint := range endpoints {
			id := strconv.Itoa(int(endpoint.Id))

			nodeNode := &domain.ResNode{Name: endpoint.Name + "(节点)", Type: _const.ResNode,
				Id: id, HostId: hostNode.Id, Key: string(_const.ResNode) + "-" + id}
			hostNode.Children = append(hostNode.Children, nodeNode)

			containerFolderNode := &domain.ResNode{Name: "实例", Type: _const.ResFolder,
				Id: id + "-folder-vms", Key: id + "-folder-container"}
			nodeNode.Children = append(nodeNode.Children, containerFolderNode)

			imageFolderNode := &domain.ResNode{Name: "镜像", Type: _const.ResFolder,
				Id: id + "-folder-templs", Key: id + "-folder-image"}
			nodeNode.Children = append(nodeNode.Children, imageFolderNode)

			containers, _ := portainer.ListContainers(endpoint.Id)
			for _, container := range containers {
				containerId := container.ID
				name := getContainerName(strings.Join(container.Names, "/"))

				vmNode := &domain.ResNode{Name: name, Type: _const.ResContainer, IsTemplate: false,
					Id: container.ID, HostId: hostNode.Id, NodeId: nodeNode.Id,
					Key: string(_const.ResContainer) + "-" + containerId}
				containerFolderNode.Children = append(containerFolderNode.Children, vmNode)
			}

			images, _ := portainer.ListImages(endpoint.Id)
			for _, image := range images {
				containerId := image.ID

				path := ""
				if len(image.RepoTags) > 0 {
					path = strings.Join(image.RepoTags, "/")
				} else if len(image.RepoDigests) > 0 {
					path = strings.Join(image.RepoDigests, "/")
				}
				name := getImageName(path)

				vmNode := &domain.ResNode{Name: name, Path: path, Type: _const.ResImage, IsTemplate: false,
					Id: image.ID, HostId: hostNode.Id, NodeId: nodeNode.Id,
					Key: string(_const.ResContainer) + "-" + containerId}
				imageFolderNode.Children = append(imageFolderNode.Children, vmNode)
			}
		}
	}

	return
}

func getContainerName(path string) string {
	if string(path[0]) == "/" {
		return path[1:]
	}
	return path
}

func getImageName(path string) string {
	arr := strings.Split(path, "/")
	if len(arr) <= 2 {
		return path
	}

	name := strings.Join(arr[len(arr)-2:], "/")
	return name
}

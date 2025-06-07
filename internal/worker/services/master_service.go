package services

import (
	"fmt"

	"github.com/robatipoor/task-scheduler/internal/worker/client"
	"github.com/robatipoor/task-scheduler/internal/worker/config"
)

type MasterService struct {
	config       config.Configure
	masterClient client.MasterClientInterface
}

type MasterServiceInterface interface {
	Register() error
}

func NewMasterService(config config.Configure, masterClient client.MasterClientInterface) *MasterService {
	return &MasterService{config: config, masterClient: masterClient}
}

func (ms *MasterService) Register() error {
	baseUrl := fmt.Sprintf("%s%s:%s", ms.config.Server.Schema, ms.config.Server.Host, ms.config.Server.Port)
	statusCode, err := ms.masterClient.Register(ms.config.Master.Url, client.WorkerRegisterRequest{BaseUrl: baseUrl})
	if err != nil {
		return err
	}

	if *statusCode != 200 {
		return fmt.Errorf("Failed register")
	}

	return nil
}

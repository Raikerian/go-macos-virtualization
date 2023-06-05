package macos

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/Code-Hex/vz/v3"
	"github.com/raikerian/go-macos-virtualization/pkg/utils"
)

type Manager struct {
	vm *vz.VirtualMachine
}

func NewManager(cpuCount uint, memorySize uint64) (*Manager, error) {
	platformConfig, err := utils.SetupMacPlatformConfiguration()
	if err != nil {
		return nil, err
	}
	config, err := utils.CreateVMConfiguration(platformConfig, cpuCount, memorySize)
	if err != nil {
		return nil, err
	}
	vm, err := vz.NewVirtualMachine(config)
	if err != nil {
		return nil, err
	}

	return &Manager{
		vm: vm,
	}, nil
}

func (m *Manager) Run(ctx context.Context) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := m.vm.Start(); err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		for {
			select {
			case newState := <-m.vm.StateChangedNotify():
				if newState == vz.VirtualMachineStateRunning {
					log.Println("VM is running")
				}
				if newState == vz.VirtualMachineStateStopped || newState == vz.VirtualMachineStateStopping {
					log.Println("stopped state")
					errCh <- nil
					return
				}
			case err := <-errCh:
				errCh <- fmt.Errorf("failed to start vm: %w", err)
				return
			}
		}
	}()

	// cleanup is this function is useful when finished graphic application.
	cleanup := func() {
		for i := 1; m.vm.CanRequestStop(); i++ {
			result, err := m.vm.RequestStop()
			log.Printf("sent stop request(%d): %t, %v", i, result, err)
			time.Sleep(time.Second * 3)
			if i > 3 {
				log.Println("call stop")
				if err := m.vm.Stop(); err != nil {
					log.Println("stop with error", err)
				}
			}
		}
		log.Println("finished cleanup")
	}

	runtime.LockOSThread()
	// TODO: this can be headlesss in the future, give a param
	m.vm.StartGraphicApplication(960, 600)
	runtime.UnlockOSThread()

	cleanup()

	return <-errCh
}

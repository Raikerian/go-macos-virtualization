package macos

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/Code-Hex/vz/v3"
	"github.com/raikerian/go-macos-virtualization/pkg/utils"
)

const restoreImageMaxRetry = 3

func Install(ctx context.Context) error {
	if err := utils.CreateVMBundle(); err != nil {
		return fmt.Errorf("failed to VM.bundle in home directory: %w", err)
	}

	restoreImagePath := utils.GetRestoreImagePath()
	var restoreImage *vz.MacOSRestoreImage
	var err error
	for i := 0; i < restoreImageMaxRetry; i++ {
		if _, err := os.Stat(restoreImagePath); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			if err := downloadRestoreImage(ctx, restoreImagePath); err != nil {
				if errRemove := os.Remove(restoreImagePath); errRemove != nil {
					return fmt.Errorf("failed to delete restore image: %w", errRemove)
				}
				continue
			}
		}
		restoreImage, err = vz.LoadMacOSRestoreImageFromPath(restoreImagePath)
		if err != nil {
			if errRemove := os.Remove(restoreImagePath); errRemove != nil {
				return fmt.Errorf("failed to delete restore image: %w", errRemove)
			}
			continue
		}

		// If everything goes well, break the loop
		break
	}

	if _, err := os.Stat(restoreImagePath); os.IsNotExist(err) {
		return fmt.Errorf("failed to download and load restore image after 3 attempts")
	}

	configurationRequirements := restoreImage.MostFeaturefulSupportedConfiguration()
	config, err := setupVirtualMachineWithMacOSConfigurationRequirements(
		configurationRequirements,
	)
	if err != nil {
		return fmt.Errorf("failed to setup config: %w", err)
	}

	vm, err := vz.NewVirtualMachine(config)
	if err != nil {
		return err
	}

	installer, err := vz.NewMacOSInstaller(vm, restoreImagePath)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("install has been cancelled")
				return
			case <-installer.Done():
				fmt.Println("install has been completed")
				return
			case <-ticker.C:
				fmt.Printf("install: %.3f%%\r", installer.FractionCompleted()*100)
			}
		}
	}()

	return installer.Install(ctx)
}

func downloadRestoreImage(ctx context.Context, destPath string) error {
	progress, err := vz.FetchLatestSupportedMacOSRestoreImage(ctx, destPath)
	if err != nil {
		return err
	}

	fmt.Printf("download restore image in %q\n", destPath)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("download has been cancelled")
			return ctx.Err()
		case <-progress.Finished():
			fmt.Println("download has been completed")
			return progress.Err()
		case <-ticker.C:
			fmt.Printf("download: %.3f%%\n", progress.FractionCompleted()*100)
		}
	}
}

func setupVirtualMachineWithMacOSConfigurationRequirements(macOSConfiguration *vz.MacOSConfigurationRequirements) (*vz.VirtualMachineConfiguration, error) {
	platformConfig, err := utils.CreateMacPlatformConfiguration(macOSConfiguration)
	if err != nil {
		return nil, fmt.Errorf("failed to create mac platform config: %w", err)
	}
	// we utilize a full power for installer to make it smooth
	return utils.CreateVMConfiguration(platformConfig, computeCPUCount(), computeMemorySize())
}

func computeCPUCount() uint {
	totalAvailableCPUs := runtime.NumCPU()
	virtualCPUCount := uint(totalAvailableCPUs - 1)
	if virtualCPUCount <= 1 {
		virtualCPUCount = 1
	}
	maxAllowed := vz.VirtualMachineConfigurationMaximumAllowedCPUCount()
	if virtualCPUCount > maxAllowed {
		virtualCPUCount = maxAllowed
	}
	minAllowed := vz.VirtualMachineConfigurationMinimumAllowedCPUCount()
	if virtualCPUCount < minAllowed {
		virtualCPUCount = minAllowed
	}
	return virtualCPUCount
}

func computeMemorySize() uint64 {
	// We arbitrarily choose 4GB.
	memorySize := uint64(4 * 1024 * 1024 * 1024)
	maxAllowed := vz.VirtualMachineConfigurationMaximumAllowedMemorySize()
	if memorySize > maxAllowed {
		memorySize = maxAllowed
	}
	minAllowed := vz.VirtualMachineConfigurationMinimumAllowedMemorySize()
	if memorySize < minAllowed {
		memorySize = minAllowed
	}
	return memorySize
}

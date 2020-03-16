package purge

import (
	"fmt"
	"path/filepath"

	"github.com/flant/werf/pkg/container_runtime"

	"github.com/spf13/cobra"

	"github.com/flant/logboek"
	"github.com/flant/shluz"

	"github.com/flant/werf/cmd/werf/common"
	"github.com/flant/werf/pkg/cleaning"
	"github.com/flant/werf/pkg/docker"
	"github.com/flant/werf/pkg/docker_registry"
	"github.com/flant/werf/pkg/storage"
	"github.com/flant/werf/pkg/werf"
)

var commonCmdData common.CmdData

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "purge",
		DisableFlagsInUseLine: true,
		Short:                 "Purge project images from images repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := common.ProcessLogOptions(&commonCmdData); err != nil {
				common.PrintHelp(cmd)
				return err
			}
			common.LogVersion()

			return common.LogRunningTime(func() error {
				return runPurge()
			})
		},
	}

	common.SetupDir(&commonCmdData, cmd)
	common.SetupTmpDir(&commonCmdData, cmd)
	common.SetupHomeDir(&commonCmdData, cmd)

	common.SetupStagesStorage(&commonCmdData, cmd)
	common.SetupImagesRepo(&commonCmdData, cmd)
	common.SetupImagesRepoMode(&commonCmdData, cmd)
	common.SetupDockerConfig(&commonCmdData, cmd, "Command needs granted permissions to delete images from the specified images repo")
	common.SetupInsecureRegistry(&commonCmdData, cmd)
	common.SetupSkipTlsVerifyRegistry(&commonCmdData, cmd)

	common.SetupLogOptions(&commonCmdData, cmd)
	common.SetupLogProjectDir(&commonCmdData, cmd)

	common.SetupDryRun(&commonCmdData, cmd)

	return cmd
}

func runPurge() error {
	if err := werf.Init(*commonCmdData.TmpDir, *commonCmdData.HomeDir); err != nil {
		return fmt.Errorf("initialization error: %s", err)
	}

	if err := shluz.Init(filepath.Join(werf.GetServiceDir(), "locks")); err != nil {
		return err
	}

	if err := docker_registry.Init(docker_registry.APIOptions{
		InsecureRegistry:      *commonCmdData.InsecureRegistry,
		SkipTlsVerifyRegistry: *commonCmdData.SkipTlsVerifyRegistry,
	}); err != nil {
		return err
	}

	if err := docker.Init(*commonCmdData.DockerConfig, *commonCmdData.LogVerbose, *commonCmdData.LogDebug); err != nil {
		return err
	}

	projectDir, err := common.GetProjectDir(&commonCmdData)
	if err != nil {
		return fmt.Errorf("getting project dir failed: %s", err)
	}

	common.ProcessLogProjectDir(&commonCmdData, projectDir)

	werfConfig, err := common.GetRequiredWerfConfig(projectDir, true)
	if err != nil {
		return fmt.Errorf("unable to load werf config: %s", err)
	}

	logboek.LogOptionalLn()

	projectName := werfConfig.Meta.Project

	stagesStorageAddress, err := common.GetStagesStorageAddress(&commonCmdData)
	if err != nil {
		return err
	}
	containerRuntime := &container_runtime.LocalDockerServerRuntime{}
	stagesStorage, err := storage.NewStagesStorage(
		stagesStorageAddress,
		containerRuntime,
		docker_registry.APIOptions{
			InsecureRegistry:      *commonCmdData.InsecureRegistry,
			SkipTlsVerifyRegistry: *commonCmdData.SkipTlsVerifyRegistry,
		},
	)
	if err != nil {
		return err
	}

	stagesStorageCache := storage.NewFileStagesStorageCache(filepath.Join(werf.GetLocalCacheDir(), "stages_storage"))
	_ = stagesStorageCache // FIXME

	storageLockManager := &storage.FileLockManager{}
	_ = storageLockManager // FIXME

	imagesRepoAddress, err := common.GetImagesRepoAddress(projectName, &commonCmdData)
	if err != nil {
		return err
	}
	imagesRepoMode, err := common.GetImagesRepoMode(&commonCmdData)
	if err != nil {
		return err
	}
	imagesRepoManager, err := storage.GetImagesRepoManager(imagesRepoAddress, imagesRepoMode)
	if err != nil {
		return err
	}
	imagesRepo, err := storage.NewImagesRepo(
		projectName,
		imagesRepoManager,
		docker_registry.APIOptions{
			InsecureRegistry:      *commonCmdData.InsecureRegistry,
			SkipTlsVerifyRegistry: *commonCmdData.SkipTlsVerifyRegistry,
		},
	)
	if err != nil {
		return err
	}

	imageNameList, err := common.GetManagedImagesNames(projectName, stagesStorage, werfConfig)
	if err != nil {
		return err
	}
	logboek.Debug.LogF("Managed images names: %v\n", imageNameList)

	imagesPurgeOptions := cleaning.ImagesPurgeOptions{
		ImageNameList: imageNameList,
		DryRun:        *commonCmdData.DryRun,
	}

	logboek.LogOptionalLn()
	if err := cleaning.ImagesPurge(imagesRepo, imagesPurgeOptions); err != nil {
		return err
	}

	return nil
}

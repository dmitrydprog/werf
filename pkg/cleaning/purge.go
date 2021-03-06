package cleaning

import (
	"context"
	"fmt"

	"github.com/werf/logboek"

	"github.com/werf/werf/pkg/image"
	"github.com/werf/werf/pkg/storage"
	"github.com/werf/werf/pkg/storage/manager"
)

type PurgeOptions struct {
	RmContainersThatUseWerfImages bool
	DryRun                        bool
}

func Purge(ctx context.Context, projectName string, storageManager *manager.StorageManager, storageLockManager storage.LockManager, options PurgeOptions) error {
	m := newPurgeManager(projectName, storageManager, options)

	if lock, err := storageLockManager.LockStagesAndImages(ctx, projectName, storage.LockStagesAndImagesOptions{GetOrCreateImagesOnly: false}); err != nil {
		return fmt.Errorf("unable to lock stages and images: %s", err)
	} else {
		defer storageLockManager.Unlock(ctx, lock)
	}

	return m.run(ctx)
}

func newPurgeManager(projectName string, storageManager *manager.StorageManager, options PurgeOptions) *purgeManager {
	return &purgeManager{
		StorageManager:                storageManager,
		ProjectName:                   projectName,
		RmContainersThatUseWerfImages: options.RmContainersThatUseWerfImages,
		DryRun:                        options.DryRun,
	}
}

type purgeManager struct {
	StorageManager                *manager.StorageManager
	ProjectName                   string
	RmContainersThatUseWerfImages bool
	DryRun                        bool
}

func (m *purgeManager) run(ctx context.Context) error {
	if err := logboek.Context(ctx).Default().LogProcess("Deleting stages").DoError(func() error {
		stages, err := m.StorageManager.GetStageDescriptionList(ctx)
		if err != nil {
			return err
		}

		return m.deleteStages(ctx, stages)
	}); err != nil {
		return err
	}

	if err := logboek.Context(ctx).Default().LogProcess("Deleting imports metadata").DoError(func() error {
		importMetadataIDs, err := m.StorageManager.StagesStorage.GetImportMetadataIDs(ctx, m.ProjectName)
		if err != nil {
			return err
		}

		return m.deleteImportsMetadata(ctx, importMetadataIDs)
	}); err != nil {
		return err
	}

	if err := logboek.Context(ctx).Default().LogProcess("Deleting managed images").DoError(func() error {
		managedImages, err := m.StorageManager.StagesStorage.GetManagedImages(ctx, m.ProjectName)
		if err != nil {
			return err
		}

		if err := m.deleteManagedImages(ctx, managedImages); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	if err := logboek.Context(ctx).Default().LogProcess("Deleting images metadata").DoError(func() error {
		_, imageMetadataByImageName, err := m.StorageManager.StagesStorage.GetAllAndGroupImageMetadataByImageName(ctx, m.ProjectName, []string{})
		if err != nil {
			return err
		}

		for imageNameID, stageIDCommitList := range imageMetadataByImageName {
			if err := m.deleteImageMetadata(ctx, imageNameID, stageIDCommitList); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (m *purgeManager) deleteStages(ctx context.Context, stages []*image.StageDescription) error {
	deleteStageOptions := manager.ForEachDeleteStageOptions{
		DeleteImageOptions: storage.DeleteImageOptions{
			RmiForce: true,
		},
		FilterStagesAndProcessRelatedDataOptions: storage.FilterStagesAndProcessRelatedDataOptions{
			SkipUsedImage:            false,
			RmForce:                  m.RmContainersThatUseWerfImages,
			RmContainersThatUseImage: m.RmContainersThatUseWerfImages,
		},
	}

	return deleteStages(ctx, m.StorageManager, m.DryRun, deleteStageOptions, stages)
}

func (m *purgeManager) deleteImportsMetadata(ctx context.Context, importsMetadataIDs []string) error {
	return deleteImportsMetadata(ctx, m.ProjectName, m.StorageManager, importsMetadataIDs, m.DryRun)
}

func (m *purgeManager) deleteManagedImages(ctx context.Context, managedImages []string) error {
	if m.DryRun {
		for _, managedImage := range managedImages {
			logboek.Context(ctx).Default().LogFDetails("  name: %s\n", managedImage)
			logboek.Context(ctx).LogOptionalLn()
		}
		return nil
	}

	return m.StorageManager.ForEachRmManagedImage(ctx, m.ProjectName, managedImages, func(ctx context.Context, managedImage string, err error) error {
		if err != nil {
			if err := handleDeletionError(err); err != nil {
				return err
			}

			logboek.Context(ctx).Warn().LogF("WARNING: Managed image %s deletion failed: %s\n", managedImage, err)

			return nil
		}

		logboek.Context(ctx).Default().LogFDetails("  name: %s\n", managedImage)

		return nil
	})
}

func (m *purgeManager) deleteImageMetadata(ctx context.Context, imageNameOrID string, stageIDCommitList map[string][]string) error {
	return deleteImageMetadata(ctx, m.ProjectName, m.StorageManager, imageNameOrID, stageIDCommitList, m.DryRun)
}

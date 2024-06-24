package dbtest

import (
	"github.com/Housiadas/backend-system/foundation/docker"
)

func migrateUp(dbTestURL string) (docker.Container, error) {
	dockerArgs := []string{
		"--entrypoint", "migrate",
	}
	appArgs := []string{
		"-path", "./business/data/migrations",
		"-database", dbTestURL,
		"up",
	}

	c, err := docker.StartContainer(MigrateImage, MigrateName, "", "", dockerArgs, appArgs)
	return c, err
}

//func migrateDown(dbTestURL string) (docker.Container, error) {
//	dockerArgs := []string{
//		"--entrypoint", "migrate",
//	}
//	appArgs := []string{
//		"-path", "./business/data/migrations",
//		"-database", dbTestURL,
//		"down",
//	}
//
//	c, err := docker.StartContainer(MigrateImage, MigrateName, "", dockerArgs, appArgs)
//	return c, err
//}

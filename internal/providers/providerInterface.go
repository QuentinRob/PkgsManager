package providers

import (
	"qrobcis/pkgsmanager/internal/models"
)

type PackageProvider interface {
	InstallPackage(pkgConfiguration models.PackageConfiguration) (err error, cmdErr error)
	UpdateRegistry() (err error, cmdErr error)
	CleanRegistry() (err error, cmdErr error)
}

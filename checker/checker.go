package checker

import "github.com/nofuture17/gonkey/models"

type CheckerInterface interface {
	Check(models.TestInterface, *models.Result) ([]error, error)
}

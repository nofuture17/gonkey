package output

import (
	"github.com/nofuture17/gonkey/models"
)

type OutputInterface interface {
	Process(models.TestInterface, *models.Result) error
}

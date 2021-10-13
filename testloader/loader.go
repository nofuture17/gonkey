package testloader

import (
	"github.com/nofuture17/gonkey/models"
)

type LoaderInterface interface {
	Load() (chan models.TestInterface, error)
}

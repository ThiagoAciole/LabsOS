package providers

import (
	"labsos/backend/internal/platform"
	"labsos/backend/providers/linux"
)

func New() platform.Provider {
	return linux.New()
}

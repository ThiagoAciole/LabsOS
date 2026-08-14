package providers

import (
	"fmt"

	"labsos/backend/internal/platform"
	"labsos/backend/providers/linux"
	"labsos/backend/providers/mock"
)

func New(mode string) (platform.Provider, error) {
	switch mode {
	case "", "mock":
		return mock.New(), nil
	case "linux":
		return linux.New(), nil
	default:
		return nil, fmt.Errorf("invalid LABSOS_MODE %q", mode)
	}
}

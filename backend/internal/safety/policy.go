package safety

import (
	"net"
	"os"
	"strings"
)

// Development defaults are loopback-only and simulated until explicitly enabled.
func Address(requested string) string {
	if strings.EqualFold(os.Getenv("LABSOS_ALLOW_REMOTE"), "true") {
		return requested
	}
	_, port, err := net.SplitHostPort(requested)
	if err != nil || port == "" {
		port = "8080"
	}
	return net.JoinHostPort("127.0.0.1", port)
}

func RealOperationsEnabled() bool {
	return strings.EqualFold(os.Getenv("LABSOS_ENABLE_REAL_OPERATIONS"), "true") && strings.TrimSpace(os.Getenv("LABSOS_CONFIRM_REAL_OPERATIONS")) != ""
}

func NetworkChangesEnabled() bool {
	return strings.EqualFold(os.Getenv("LABSOS_ENABLE_NETWORK_CHANGES"), "true") && RealOperationsEnabled()
}

package api

import "testing"

func TestValidInterfaceName(t *testing.T) {
	for _, value := range []string{"eth0", "enp3s0", "wlan_0"} {
		if !validInterfaceName(value) {
			t.Errorf("expected valid interface: %s", value)
		}
	}
	for _, value := range []string{"", "eth0;rm", "../../x", "this-interface-name-is-too-long"} {
		if validInterfaceName(value) {
			t.Errorf("expected invalid interface: %s", value)
		}
	}
}

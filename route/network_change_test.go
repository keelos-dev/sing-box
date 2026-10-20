package route

import (
	"net"
	"net/netip"
	"strings"
	"testing"

	"github.com/sagernet/sing/common/control"
)

func TestDescribeInterfaceChange(t *testing.T) {
	base := control.Interface{
		Index: 11, Name: "Ethernet", MTU: 1500, Flags: net.FlagUp,
		Addresses: []netip.Prefix{netip.MustParsePrefix("192.168.11.226/24"), netip.MustParsePrefix("240f::1/64")},
	}
	old := &interfaceSnapshot{iface: base, dnsServers: []netip.Addr{netip.MustParseAddr("8.8.8.8")}}

	same := *old
	if got := describeInterfaceChange(old, &same); got != "no visible change" {
		t.Fatal(got)
	}

	moved := base
	moved.Addresses = []netip.Prefix{netip.MustParsePrefix("192.168.11.226/24"), netip.MustParsePrefix("240f::2/64")}
	got := describeInterfaceChange(old, &interfaceSnapshot{iface: moved, dnsServers: []netip.Addr{netip.MustParseAddr("fe80::1"), netip.MustParseAddr("8.8.8.8")}})
	for _, want := range []string{"addresses +[240f::2/64] -[240f::1/64]", "dns servers [8.8.8.8] -> [fe80::1 8.8.8.8]"} {
		if !strings.Contains(got, want) {
			t.Fatalf("%q missing %q", got, want)
		}
	}

	vpn := base
	vpn.Index, vpn.Name = 29, "VPN"
	if got := describeInterfaceChange(old, &interfaceSnapshot{iface: vpn, dnsServers: old.dnsServers}); !strings.Contains(got, "interface Ethernet (11) -> VPN (29)") {
		t.Fatal(got)
	}
	if describeInterfaceChange(nil, old) != "initial" || describeInterfaceChange(old, nil) != "default interface lost" {
		t.Fatal("initial / lost")
	}
}

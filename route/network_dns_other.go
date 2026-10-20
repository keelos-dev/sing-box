//go:build !windows

package route

import "net/netip"

func interfaceDNSServers(int) []netip.Addr { return nil }

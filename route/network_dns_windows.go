package route

import (
	"net/netip"
	"unsafe"

	"golang.org/x/sys/windows"
)

// interfaceDNSServers returns the DNS servers configured on the interface, in the order Windows reports them.
func interfaceDNSServers(index int) []netip.Addr {
	size := uint32(15 * 1024)
	for range 3 {
		buffer := make([]byte, size)
		first := (*windows.IpAdapterAddresses)(unsafe.Pointer(&buffer[0]))
		err := windows.GetAdaptersAddresses(windows.AF_UNSPEC,
			windows.GAA_FLAG_SKIP_ANYCAST|windows.GAA_FLAG_SKIP_MULTICAST, 0, first, &size)
		if err == windows.ERROR_BUFFER_OVERFLOW {
			continue
		}
		if err != nil {
			return nil
		}
		for adapter := first; adapter != nil; adapter = adapter.Next {
			if int(adapter.IfIndex) != index && int(adapter.Ipv6IfIndex) != index {
				continue
			}
			var servers []netip.Addr
			for dns := adapter.FirstDnsServerAddress; dns != nil; dns = dns.Next {
				if addr, ok := netip.AddrFromSlice(dns.Address.IP()); ok {
					servers = append(servers, addr.Unmap())
				}
			}
			return servers
		}
		return nil
	}
	return nil
}

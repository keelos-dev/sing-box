package route

import (
	"bytes"
	"net/netip"
	"slices"
	"strings"

	"github.com/sagernet/sing/common/control"
	F "github.com/sagernet/sing/common/format"
)

// interfaceSnapshot is what the default interface monitor compares between two updates:
// the interface itself plus its DNS servers.
type interfaceSnapshot struct {
	iface      control.Interface
	dnsServers []netip.Addr
}

func takeInterfaceSnapshot(iface *control.Interface) *interfaceSnapshot {
	if iface == nil {
		return nil
	}
	return &interfaceSnapshot{
		iface:      *iface,
		dnsServers: interfaceDNSServers(iface.Index),
	}
}

// describeInterfaceChange explains why the default interface monitor reported an update,
// so the log shows which field changed instead of only "updated default interface".
func describeInterfaceChange(old, current *interfaceSnapshot) string {
	if current == nil {
		return "default interface lost"
	}
	if old == nil {
		return "initial"
	}
	var changes []string
	o, n := old.iface, current.iface
	if o.Index != n.Index || o.Name != n.Name {
		changes = append(changes, F.ToString("interface ", o.Name, " (", o.Index, ") -> ", n.Name, " (", n.Index, ")"))
	}
	if o.MTU != n.MTU {
		changes = append(changes, F.ToString("mtu ", o.MTU, " -> ", n.MTU))
	}
	if !bytes.Equal(o.HardwareAddr, n.HardwareAddr) {
		changes = append(changes, F.ToString("hardware address ", o.HardwareAddr, " -> ", n.HardwareAddr))
	}
	if o.Flags != n.Flags {
		changes = append(changes, F.ToString("flags ", o.Flags, " -> ", n.Flags))
	}
	if added, removed := diffPrefixes(o.Addresses, n.Addresses); len(added)+len(removed) > 0 {
		changes = append(changes, F.ToString("addresses +", joinPrefixes(added), " -", joinPrefixes(removed)))
	}
	if !slices.Equal(old.dnsServers, current.dnsServers) {
		changes = append(changes, F.ToString("dns servers ", joinAddrs(old.dnsServers), " -> ", joinAddrs(current.dnsServers)))
	}
	if len(changes) == 0 {
		return "no visible change"
	}
	return strings.Join(changes, "; ")
}

func diffPrefixes(old, current []netip.Prefix) (added, removed []netip.Prefix) {
	for _, p := range current {
		if !slices.Contains(old, p) {
			added = append(added, p)
		}
	}
	for _, p := range old {
		if !slices.Contains(current, p) {
			removed = append(removed, p)
		}
	}
	return
}

func joinPrefixes(prefixes []netip.Prefix) string {
	parts := make([]string, len(prefixes))
	for i, p := range prefixes {
		parts[i] = p.String()
	}
	return "[" + strings.Join(parts, " ") + "]"
}

func joinAddrs(addrs []netip.Addr) string {
	parts := make([]string, len(addrs))
	for i, a := range addrs {
		parts[i] = a.String()
	}
	return "[" + strings.Join(parts, " ") + "]"
}

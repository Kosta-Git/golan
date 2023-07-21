package interfaces

import (
    "net"
    "testing"
)

func TestFindInterfaceByIp(t *testing.T) {
    t.Run("Interface exists", func(t *testing.T) {
        ifaces, _ := net.Interfaces()
        for _, iface := range ifaces {
            if iface.Flags&net.FlagUp != net.FlagUp || iface.Flags&net.FlagRunning != net.FlagRunning {
                continue
            }
            addrs, _ := iface.Addrs()
            if len(addrs) == 0 {
                continue
            }
            addr := addrs[0]
            ip, _, _ := net.ParseCIDR(addr.String())
            _, err := FindInterfaceByIp(ip)
            if err != nil {
                t.Errorf("FindInterfaceByIp(%v) error = %v", ip, err)
                return
            }
        }
    })

    t.Run("Interface doesn't exists", func(t *testing.T) {
        _, err := FindInterfaceByIp(net.ParseIP("255.255.255.255"))
        if err == nil {
            t.Errorf("FindInterfaceByIp() expected error")
            return
        }
    })
}

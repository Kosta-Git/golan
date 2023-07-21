package interfaces

import (
    "errors"
    "net"
)

func FindInterfaceByIp(ip net.IP) (net.Interface, error) {
    ifaces, err := net.Interfaces()
    if err != nil {
        return net.Interface{}, err
    }
    for _, iface := range ifaces {
        addrs, err := iface.Addrs()
        if err != nil {
            return net.Interface{}, err
        }

        for _, addr := range addrs {
            _, network, err := net.ParseCIDR(addr.String())
            if err != nil {
                return net.Interface{}, err
            }
            if network.Contains(ip) {
                return iface, nil
            }
        }
    }
    return net.Interface{}, errors.New("interface not found")
}

func FindAnyUpAndRunningInterface() (net.Interface, error) {
    ifaces, err := net.Interfaces()
    if err != nil {
        return net.Interface{}, err
    }
    for _, iface := range ifaces {
        if iface.Flags&net.FlagUp == net.FlagUp && iface.Flags&net.FlagRunning == net.FlagRunning {
            return iface, nil
        }
    }
    return net.Interface{}, errors.New("no interface found")
}

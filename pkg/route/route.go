package route

import (
    "errors"
    "golan/pkg/interfaces"
    "golan/pkg/mask"
    "net"
    "os/exec"
    "strconv"
)

type Route struct {
    destination  net.IP
    netMask      net.IPMask
    gateway      net.IP
    iface        net.Interface
    hasInterface bool
    metric       int
}

func NewRouteFindInterface(destination net.IP, netmask net.IPMask, gateway net.IP, iface net.IP, metric int) Route {
    ifaceObj := net.Interface{}
    hasIface := false
    var err error
    if iface != nil {
        ifaceObj, err = interfaces.FindInterfaceByIp(iface)
        if err == nil {
            hasIface = true
        }
    }

    return Route{
        destination:  destination,
        netMask:      netmask,
        gateway:      gateway,
        iface:        ifaceObj,
        hasInterface: hasIface,
        metric:       metric,
    }
}

func NewRoute(destination net.IP, netmask net.IPMask, gateway net.IP, iface net.Interface, hasInterface bool, metric int) Route {
    return Route{
        destination:  destination,
        netMask:      netmask,
        gateway:      gateway,
        iface:        iface,
        hasInterface: hasInterface,
        metric:       metric,
    }
}

func (r *Route) Metric() int {
    return r.metric
}

func (r *Route) Interface() net.Interface {
    return r.iface
}

func (r *Route) Gateway() net.IP {
    return r.gateway
}

func (r *Route) ExecChangeMetric(updatedMetric int) error {
    r.metric = updatedMetric
    return r.exec("CHANGE")
}

func (r *Route) ExecAddRoute() error {
    return r.exec("ADD")
}

func (r *Route) ExecDeleteRoute() error {
    return r.exec("DELETE")
}

func (r *Route) exec(operation string) error {
    if !r.hasInterface || r.gateway == nil {
        return errors.New("no interface found for route")
    }
    cmd := exec.Command(
        "route",
        operation,
        r.destination.String(),
        "MASK",
        mask.ToIp(r.netMask).String(),
        r.gateway.String(),
        "METRIC",
        strconv.Itoa(r.metric),
        "IF",
        strconv.Itoa(r.iface.Index),
    )
    return cmd.Run()
}

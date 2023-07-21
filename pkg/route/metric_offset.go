package route

import (
    "net"
)

func FindMetricOffset(iface net.Interface, gateway net.IP) (int, error) {
    route := Route{
        destination:  net.ParseIP("255.255.255.255"),
        netMask:      net.IPv4Mask(255, 255, 255, 255),
        gateway:      gateway,
        iface:        iface,
        hasInterface: true,
        metric:       1,
    }

    err := route.ExecAddRoute()
    if err != nil {
        return 0, err
    }

    foundRoutes := FindRoute(route)
    if len(foundRoutes) == 0 {
        return 0, nil
    }

    err = route.ExecDeleteRoute()
    if err != nil {
        return 0, err
    }

    return foundRoutes[0].metric - 1, nil
}

func MetricCanBeLowerThan(targetRoute *Route, interfaceOffset int) bool {
    return (targetRoute.metric - 1) > interfaceOffset
}

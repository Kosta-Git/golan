package route

import (
    "fmt"
    "golan/pkg/mask"
    "log"
    "net"
    "sort"
)

// Target is a struct that represents a minimal target route
// that can be used to verify if a route exists, if it doesn't
// it can create the route, if the route exists but does not match
// the requirements it can update the route or other routes to match
type Target struct {
    destination           net.IP
    minNetMask            net.IPMask
    destinationMinNetwork *net.IPNet
    gateway               net.IP
    iface                 net.Interface
    ifaceOffset           int
}

func NewTarget(route Route, interfaceMetricOffset int) Target {
    _, maskSize := route.netMask.Size()
    _, net, err := net.ParseCIDR(fmt.Sprintf("%s/%d", route.destination, maskSize))
    if err != nil {
        log.Fatal(err)
    }

    return Target{
        destination:           route.destination,
        minNetMask:            route.netMask,
        destinationMinNetwork: net,
        gateway:               route.gateway,
        iface:                 route.iface,
        ifaceOffset:           interfaceMetricOffset,
    }
}

func (t *Target) Verify(routes []Route) (bool, error) {
    routesForTarget := t.filterRoutesByInterface(routes)

    // We have no route for our gateway, we need to create one
    if len(routesForTarget) == 0 {
        return t.createRoute()
    }

    // Sort by metric ascending, we can check if our target route is the first one
    sort.Slice(routesForTarget, func(i, j int) bool {
        return routesForTarget[i].Metric() < routesForTarget[j].Metric()
    })

    found := false
    foundAtIndex := -1
    for i, route := range routesForTarget {
        if t.isMatchingRoute(route) {
            found = true
            foundAtIndex = i
            break
        }
    }

    if !found {
        return t.createRoute()
    }

    if foundAtIndex == 0 {
        log.Printf("Route for gateway %s found, and doesn't need updating.", t.gateway)
        return true, nil
    }

    targetRoute := routesForTarget[foundAtIndex]
    firstRoute := routesForTarget[0]

    if targetRoute.metric == firstRoute.metric {
        return true, nil
    }

    log.Printf("Route for gateway %s found, but needs updating.", t.gateway)

    // This is the easy case, we can just set our route before the first route
    if MetricCanBeLowerThan(&firstRoute, t.ifaceOffset) {
        log.Printf("Metric for route %v can be lower than %d, updating target route to be lower", firstRoute, t.ifaceOffset)
        err := targetRoute.ExecChangeMetric(targetRoute.metric - t.ifaceOffset - 1)
        if err != nil {
            return false, err
        }
        return true, nil
    }

    // This is the hard case, we need to update our route to be first, then add one to the second route
    // We then need to check if this affects the existing order, if yes we need to fix it
    log.Printf("Metric for route %v cannot be lower than %d, updating all routes to ensure same order", firstRoute, t.ifaceOffset)

    metricChangeMap := make(map[int]int)
    metricChangeMap[foundAtIndex] = firstRoute.metric - t.ifaceOffset
    metricChangeMap[0] = firstRoute.metric - t.ifaceOffset + 1

    for i := 1; i < len(routesForTarget); i++ {
        if i == foundAtIndex {
            continue
        }

        route := routesForTarget[i]
        if route.metric+1 == routesForTarget[i-1].metric {
            metricChangeMap[i] = route.metric - t.ifaceOffset + 1
        }
    }

    log.Printf("Preparing to apply %d modifications", len(metricChangeMap))

    for i, metric := range metricChangeMap {
        err := routesForTarget[i].ExecChangeMetric(metric)
        if err != nil {
            log.Printf("Error while updating route %+v", routesForTarget[i])
        }
    }

    return true, nil
}

func (t *Target) createRoute() (bool, error) {
    log.Printf("No route for gateway %s found, creating one", t.gateway)
    routeToCreate := NewRoute(t.destination, t.minNetMask, t.gateway, t.iface, true, 1)
    err := routeToCreate.ExecAddRoute()
    if err != nil {
        return false, err
    }
    return true, nil
}

func (t *Target) isMatchingRoute(route Route) bool {
    return t.destinationMinNetwork.Contains(route.destination) &&
        route.Gateway().Equal(t.gateway) &&
        route.Interface().Index == t.iface.Index &&
        mask.IsSmallerOrEqualTo(route.netMask, t.minNetMask)
}

func (t *Target) filterRoutesByInterface(routes []Route) []Route {
    var filteredRoutes []Route
    for _, route := range routes {
        if route.iface.Index != t.iface.Index {
            continue
        }
        filteredRoutes = append(filteredRoutes, route)
    }
    return filteredRoutes
}

package core

import (
    r "golan/pkg/route"
    "log"
)

func Init() []r.Target {
    var targets []r.Target
    for _, route := range r.RetrieveCurrentRoutesV4() {
        if route.Gateway() != nil {
            offset, err := r.FindMetricOffset(route.Interface(), route.Gateway())
            if err != nil {
                log.Fatalf("Failed to find metric offset for interface %s and gateway %s: %s", route.Interface().Name, route.Gateway(), err)
            }
            targets = append(targets, r.NewTarget(route, offset))
        }
    }
    return targets
}

package mask

import "net"

func IsSmallerOrEqualTo(toCheck net.IPMask, target net.IPMask) bool {
    if toCheck == nil || target == nil {
        return false
    }
    return ToUint32(toCheck) <= ToUint32(target)
}

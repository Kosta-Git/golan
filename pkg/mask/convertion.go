package mask

import "net"

func FromIp(mask net.IP) net.IPMask {
    if mask == nil {
        return nil
    }
    mask = mask.To4()
    return net.IPv4Mask(mask[0], mask[1], mask[2], mask[3])
}

func ToIp(mask net.IPMask) net.IP {
    if mask == nil {
        return nil
    }
    return net.IPv4(mask[0], mask[1], mask[2], mask[3])
}

func ToUint32(mask net.IPMask) uint32 {
    if mask == nil {
        return 0
    }
    return uint32(mask[0])<<24 | uint32(mask[1])<<16 | uint32(mask[2])<<8 | uint32(mask[3])
}

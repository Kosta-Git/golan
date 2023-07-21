package mask

import (
    "net"
    "reflect"
    "testing"
)

func TestFromIp(t *testing.T) {
    type args struct {
        mask net.IP
    }
    tests := []struct {
        name string
        args args
        want net.IPMask
    }{
        {
            name: "/32",
            args: args{
                mask: net.IPv4(255, 255, 255, 255),
            },
            want: net.IPv4Mask(255, 255, 255, 255),
        },
        {
            name: "/24",
            args: args{
                mask: net.IPv4(255, 255, 255, 0),
            },
            want: net.IPv4Mask(255, 255, 255, 0),
        },
        {
            name: "/16",
            args: args{
                mask: net.IPv4(255, 255, 0, 0),
            },
            want: net.IPv4Mask(255, 255, 0, 0),
        },
        {
            name: "/8",
            args: args{
                mask: net.IPv4(255, 0, 0, 0),
            },
            want: net.IPv4Mask(255, 0, 0, 0),
        },
        {
            name: "/0",
            args: args{
                mask: net.IPv4(0, 0, 0, 0),
            },
            want: net.IPv4Mask(0, 0, 0, 0),
        },
        {
            name: "/31",
            args: args{
                mask: net.IPv4(255, 255, 255, 254),
            },
            want: net.IPv4Mask(255, 255, 255, 254),
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := FromIp(tt.args.mask); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("FromIp() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestToIp(t *testing.T) {
    type args struct {
        mask net.IPMask
    }
    tests := []struct {
        name string
        args args
        want net.IP
    }{
        {
            name: "/32",
            args: args{
                mask: net.IPv4Mask(255, 255, 255, 255),
            },
            want: net.ParseIP("255.255.255.255"),
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := ToIp(tt.args.mask); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ToIp() = %v, want %v", got, tt.want)
            }
        })
    }
}

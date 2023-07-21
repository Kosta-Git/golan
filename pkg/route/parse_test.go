package route

import (
    "bufio"
    "net"
    "reflect"
    "strings"
    "testing"
)

func Test_parse(t *testing.T) {
    type args struct {
        scanner *bufio.Scanner
    }
    tests := []struct {
        name string
        args args
        want []Route
    }{
        {
            name: "parseRoutePrint4",
            args: args{
                scanner: bufio.NewScanner(strings.NewReader(`
===========================================================================
Interface List
  3...aa aa aa aa aa aa ......Realtek PCIe GbE Family Controller
===========================================================================

IPv4 Route Table
===========================================================================
Active Routes:
Network Destination        Netmask          Gateway       Interface  Metric
          0.0.0.0          0.0.0.0      192.168.1.1     255.168.1.32     25
        127.0.0.0        255.0.0.0         On-link         255.0.0.1    331
      169.254.0.0      255.255.0.0         On-link      255.254.29.9    291
===========================================================================
Persistent Routes:
  None
`)),
            },
            want: []Route{
                {
                    destination:  net.IPv4zero,
                    netMask:      net.IPv4Mask(0, 0, 0, 0),
                    gateway:      net.ParseIP("192.168.1.1"),
                    iface:        net.Interface{},
                    hasInterface: false,
                    metric:       25,
                },
                {
                    destination:  net.IPv4(127, 0, 0, 0),
                    netMask:      net.IPv4Mask(255, 0, 0, 0),
                    gateway:      nil,
                    iface:        net.Interface{},
                    hasInterface: false,
                    metric:       331,
                },
                {
                    destination:  net.IPv4(169, 254, 0, 0),
                    netMask:      net.IPv4Mask(255, 255, 0, 0),
                    gateway:      nil,
                    iface:        net.Interface{},
                    hasInterface: false,
                    metric:       291,
                },
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := parse(tt.args.scanner); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("parse() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_parseRoute(t *testing.T) {
    type args struct {
        line string
    }
    tests := []struct {
        name string
        args args
        want Route
    }{
        {
            name: "parseValidRoute",
            args: args{
                line: "          0.0.0.0          0.0.0.0      192.168.1.1     255.168.1.9     25",
            },
            want: Route{
                destination:  net.IPv4zero,
                netMask:      net.IPv4Mask(0, 0, 0, 0),
                gateway:      net.ParseIP("192.168.1.1"),
                iface:        net.Interface{},
                hasInterface: false,
                metric:       25,
            },
        },
        {
            name: "parseValidRouteWithOnLinkGateway",
            args: args{
                line: "          0.0.0.0          0.0.0.0      On-Link     255.168.1.9     331",
            },
            want: Route{
                destination:  net.IPv4zero,
                netMask:      net.IPv4Mask(0, 0, 0, 0),
                gateway:      nil,
                iface:        net.Interface{},
                hasInterface: false,
                metric:       25,
            },
        },
        {
            name: "parseValidRouteWith/24Mask",
            args: args{
                line: "          0.0.0.0          255.255.255.0      On-Link     255.168.1.9     331",
            },
            want: Route{
                destination:  net.IPv4zero,
                netMask:      net.IPv4Mask(255, 255, 255, 0),
                gateway:      nil,
                iface:        net.Interface{},
                hasInterface: false,
                metric:       25,
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := parseRoute(tt.args.line); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("parseRoute() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_parseRoutes(t *testing.T) {
    type args struct {
        scanner *bufio.Scanner
    }
    tests := []struct {
        name string
        args args
        want []Route
    }{
        {
            name: "parseRoutePrint4",
            args: args{
                scanner: bufio.NewScanner(strings.NewReader(`
          0.0.0.0          0.0.0.0      192.168.1.1     255.168.1.32     25
        127.0.0.0        255.0.0.0         On-link         255.0.0.1    331
      169.254.0.0      255.255.0.0         On-link      255.254.29.9    291
===========================================================================
Persistent Routes:
  None
`)),
            },
            want: []Route{
                {
                    destination:  net.IPv4zero,
                    netMask:      net.IPv4Mask(0, 0, 0, 0),
                    gateway:      net.ParseIP("192.168.1.1"),
                    iface:        net.Interface{},
                    hasInterface: false,
                    metric:       25,
                },
                {
                    destination:  net.IPv4(127, 0, 0, 0),
                    netMask:      net.IPv4Mask(255, 0, 0, 0),
                    gateway:      nil,
                    iface:        net.Interface{},
                    hasInterface: false,
                    metric:       331,
                },
                {
                    destination:  net.IPv4(169, 254, 0, 0),
                    netMask:      net.IPv4Mask(255, 255, 0, 0),
                    gateway:      nil,
                    iface:        net.Interface{},
                    hasInterface: false,
                    metric:       291,
                },
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := parseRoutes(tt.args.scanner); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("parseRoutes() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_runRoutePrint4(t *testing.T) {
    tests := []struct {
        name string
    }{
        {
            name: "runRoutePrint4",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := runRoutePrint("print", "-4"); !bufio.NewScanner(got).Scan() {
                t.Errorf("Scanner is empty")
            }
        })
    }
}

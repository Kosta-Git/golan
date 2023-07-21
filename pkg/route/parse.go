package route

import (
    "bufio"
    "golan/pkg/mask"
    "io"
    "log"
    "net"
    "os/exec"
    "strconv"
    "strings"
)

func RetrieveCurrentRoutesV4() []Route {
    stdout := runRoutePrint("print", "-4")
    scanner := bufio.NewScanner(stdout)
    scanner.Split(bufio.ScanLines)
    return parse(scanner)
}

func FindRoute(route Route) []Route {
    stdout := runRoutePrint("print", route.destination.String(), "MASK", route.netMask.String(), route.gateway.String())
    scanner := bufio.NewScanner(stdout)
    scanner.Split(bufio.ScanLines)
    return parse(scanner)
}

func parse(scanner *bufio.Scanner) []Route {
    skipHeader(scanner)
    return parseRoutes(scanner)
}

func skipHeader(scanner *bufio.Scanner) {
    equals, skipAfterEquals := 3, 2
    for scanner.Scan() {

        line := scanner.Text()

        if equals > 0 {
            if strings.HasPrefix(line, "=") {
                equals--
            }
            continue
        }

        if skipAfterEquals > 0 {
            skipAfterEquals--

            if equals == 0 && skipAfterEquals == 0 {
                break
            }

            continue
        }
    }
}

func parseRoutes(scanner *bufio.Scanner) []Route {
    var routes []Route

    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "=") {
            break
        }

        route := parseRoute(line)
        routes = append(routes, route)
    }

    return routes
}

func parseRoute(line string) Route {
    elements := strings.Fields(line)
    if len(elements) != 5 {
        log.Fatalf("Unable to parse route: %v", line)
    }

    destination := net.ParseIP(elements[0])
    netmask := mask.FromIp(net.ParseIP(elements[1]))
    gateway := net.ParseIP(elements[2])
    iface := net.ParseIP(elements[3])
    metric, err := strconv.Atoi(elements[4])
    if err != nil {
        metric = -1
    }

    return NewRouteFindInterface(destination, netmask, gateway, iface, metric)
}

func runRoutePrint(arg ...string) io.Reader {
    cmd := exec.Command("route", arg...)
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Fatalf("Unable to obtain command stdout: %v", err)
    }
    if err := cmd.Start(); err != nil {
        log.Fatalf("Unable to execute command: %v", err)
    }
    return stdout
}

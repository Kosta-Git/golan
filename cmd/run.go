package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "golan/pkg/core"
    golanSvc "golan/pkg/service"
    "golang.org/x/sys/windows/svc"
    "golang.org/x/sys/windows/svc/mgr"
    "log"
    "time"
)

var runCmd = &cobra.Command{
    Use:   "run [duration i.e. 5s]",
    Short: "Starts the golan daemon",
    Long: `Starts the golan daemon. The daemon polls the system's routing table every [poll duration] and updates the routing table accordingly.
If no poll duration is specified, the default duration is 5 seconds.
The daemon must run before any VPN connections are established or the daemon will not be able to determine the default routes.
The daemon will only work with IPv4 routes.`,
    Run: func(cmd *cobra.Command, args []string) {
        runAsService(args)
    },
}

func init() {
    rootCmd.AddCommand(runCmd)
}

func runAsService(args []string) {
    isWindowsService, err := svc.IsWindowsService()
    if err != nil {
        log.Fatalf("Failed to determine if we are running as a service: %s", err)
    }

    if isWindowsService {
        runService()
    } else {
        duration := 5 * time.Second // Default duration
        if len(args) > 0 {
            d, err := time.ParseDuration(args[0])
            if err != nil {
                fmt.Println("Invalid duration. Using default duration.")
            } else {
                duration = d
            }
        }
        runCli(duration)
    }
}

func runService() {
    serviceManager, err := mgr.Connect()
    if err != nil {
        log.Fatalf("Failed to connect to service manager: %s", err)
    }
    defer serviceManager.Disconnect()

    service, err := serviceManager.OpenService("GolanService")
    if err != nil {
        log.Fatalf("Failed to open service: %s", err)
    }
    defer service.Close()

    err = svc.Run("GolanService", &golanSvc.GolanService{})
    if err != nil {
        log.Fatalf("Failed to run service: %s", err)
    }
}

func runCli(duration time.Duration) {
    targets := core.Init()
    for {
        err := core.Verify(targets)
        if err != nil {
            log.Println(err)
        }
        time.Sleep(duration)
    }
}

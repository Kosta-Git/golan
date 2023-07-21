package cmd

import (
    "github.com/spf13/cobra"
    "log"
    "os"
    "os/exec"
    "strings"
)

var loopDuration string

// installCmd represents the install command
var installCmd = &cobra.Command{
    Use:   "install",
    Short: "Install golan as a windows service, must be run as administrator and from C:\\Program Files\\Golan\\golan.exe move the executable there",
    Run:   installService,
}

func init() {
    rootCmd.AddCommand(installCmd)
}

func installService(cmd *cobra.Command, args []string) {
    log.Println("Installing golan as a service")
    exePath, err := getExePath()
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    arguments := []string{
        "create",
        "GolanService",
        "binPath=",
        exePath,
        "start=",
        "auto",
        "DisplayName=",
        "GolanService",
    }
    log.Println("Running: sc", strings.Join(arguments, " "))
    out, err := exec.Command("sc", arguments...).CombinedOutput()
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    log.Println(string(out))
    log.Println("Service installed successfully, if you want to start it run: golan.exe start")
}

func getExePath() (string, error) {
    exePath, err := os.Executable()
    if err != nil {
        return "", err
    }
    return exePath, nil
}

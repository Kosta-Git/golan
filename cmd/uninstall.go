package cmd

import (
    "log"
    "os/exec"

    "github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
    Use:   "uninstall",
    Short: "Uninstalls golan windows service",
    Run: func(cmd *cobra.Command, args []string) {
        log.Println("Uninstalling the service...")

        out, err := exec.Command("sc", "delete", "GolanService").CombinedOutput()
        if err != nil {
            log.Println("Error:", err)
        }
        log.Println(string(out))
    },
}

func init() {
    rootCmd.AddCommand(uninstallCmd)
}

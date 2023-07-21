package cmd

import (
    "log"
    "os/exec"

    "github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
    Use:   "stop",
    Short: "Stops the golan windows service",
    Run: func(cmd *cobra.Command, args []string) {
        log.Println("Stopping the service...")

        out, err := exec.Command("sc", "stop", "GolanService").CombinedOutput()
        if err != nil {
            log.Println("Error:", err)
        }
        log.Println(string(out))
    },
}

func init() {
    rootCmd.AddCommand(stopCmd)
}

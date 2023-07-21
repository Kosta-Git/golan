package cmd

import (
    "log"
    "os/exec"

    "github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
    Use:   "start",
    Short: "Starts the golan windows service",
    Run: func(cmd *cobra.Command, args []string) {
        log.Println("Starting the service...")

        out, err := exec.Command("sc", "start", "GolanService").CombinedOutput()
        if err != nil {
            log.Println("Error:", err)
        }
        log.Println(string(out))
    },
}

func init() {
    rootCmd.AddCommand(startCmd)
}

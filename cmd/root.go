package cmd

import (
    "github.com/spf13/cobra"
    "os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "golan",
    Short: "Daemon to prevent VPN from overriding default routes",
    Long: `Golan is a daemon that prevents VPN connections from overriding default routes.
The daemon polls the system's routing table and updates the routing table accordingly.
It must run before any VPN connections are established or the daemon will not be able to determine the default routes.`,
    Run: func(cmd *cobra.Command, args []string) {
        runAsService(args)
    },
}

func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}

func init() {
}

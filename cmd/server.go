package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"telnet/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server with optional port",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if len(args) > 0 {
			port, err := strconv.ParseUint(args[0], 10, 16)
			if err != nil {
				fmt.Printf("error parsing port from value \"%s\"\n", args[0])
				os.Exit(1)
			}
			err = server.StartAtPort(uint16(port))
		} else {
			err = server.Start()
		}

		if err != nil {
			fmt.Printf("error starting server: %s\n", err)
		}
	},
}

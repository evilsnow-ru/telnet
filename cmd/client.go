package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"telnet/client"
	"telnet/server"
)

var clientCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to remote telnet",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		host := "localhost"
		var port uint16 = server.DefaultPort

		if len(args) > 2 {
			fmt.Println("Too many arguments...")
			os.Exit(0)
		}

		switch len(args) {
		case 2:
			host = args[0]
			inputPort, err := strconv.ParseUint(args[1], 10, 16)

			if err != nil {
				fmt.Printf("error parsing port from value \"%s\"\n", args[0])
				os.Exit(1)
			}

			port = uint16(inputPort)

		case 1:
			host = args[0]
		}

		err = client.Start(host, port)

		if err != nil {
			fmt.Printf("Error running client: %s", err)
		}
	},
}

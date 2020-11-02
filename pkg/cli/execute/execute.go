package execute

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cian911/raspberry-pi-provisioner/pkg/ssh"
	"github.com/cian911/raspberry-pi-provisioner/pkg/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "execute",
		Short: "Execute command on remote hosts",
		Long:  "TODO: Write a longer message for this.",
		Run: func(cmd *cobra.Command, args []string) {
			configFile := viper.GetString("config-file")
			parsedFile, err := yaml.ParseFile(configFile)

			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			results := make(chan string, 10)
			timeout := time.After(10 * time.Second)

			masterNodes := make(map[string]string)
			workerNodes := make(map[string]string)

			for _, v := range parsedFile.MasterN.H {
				masterNodes[v.Hostname] = v.IPAddress
			}

			for _, v := range parsedFile.WorkerN.W {
				workerNodes[v.Hostname] = v.IPAddress
			}

			// I guess we want  to try and connect to the nodes here..
			for _, host := range workerNodes {
				go func(host string) {
					results <- ssh.Execute(host)
				}(host)
			}

			for i := 0; i < len(workerNodes); i++ {
				select {
				case res := <-results:
					fmt.Println(res)
				case <-timeout:
					fmt.Println("Timeout.")
					return
				}
			}
		},
	}

	c.Flags().StringP("config-file", "", "", "Pass a configuration file to syncwave to be parsed.")

	viper.BindPFlag("config-file", c.Flags().Lookup("config-file"))

	return
}
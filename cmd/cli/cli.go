package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd defines the root cli command
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kubernetes-api-reference",
		Short:         "K8s API documentation tools",
		Long:          `Tool to build documentation from OpenAPI specification of the Kubernetes API`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("go")
			return nil
		},
	}

	cmd.PersistentFlags().StringP("file", "f", "", "OpenAPI spec file")
	cmd.MarkFlagRequired("file")

	cmd.AddCommand(ResourceslistCmd())
	cmd.AddCommand(ShowTOCCmd())
	cmd.AddCommand(GVKeysMap())

	cobra.OnInitialize(initConfig)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

// Run the cli
func Run() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}

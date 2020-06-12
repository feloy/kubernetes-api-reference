package cli

import (
	"path"

	"github.com/feloy/kubernetes-api-reference/pkg/config"
	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
	"github.com/spf13/cobra"
)

// Hugo defines the `hugo` subcommand
func Hugo() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "hugo",
		Short:         "output specification for Hugo website",
		Long:          "output the specification in a format usable for a Hugo website",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			file := cmd.Flag("file").Value.String()
			spec, err := kubernetes.NewSpec(file)
			if err != nil {
				return err
			}

			configDir := cmd.Flag("config-dir").Value.String()
			toc, err := config.LoadTOC(path.Join(configDir, "toc.yaml"))
			err = toc.PopulateAssociates(spec)
			if err != nil {
				return err
			}

			toc.AddOtherResources(spec)

			outputDir := cmd.Flag("output-dir").Value.String()
			return toc.ToHugo(outputDir)
		},
	}
	cmd.Flags().StringP("config-dir", "c", "", "Directory containing documentation configuration")
	cmd.MarkFlagRequired("config-dir")
	cmd.Flags().StringP("output-dir", "o", "", "Directory to write markdown files")
	cmd.MarkFlagRequired("output-dir")

	return cmd
}

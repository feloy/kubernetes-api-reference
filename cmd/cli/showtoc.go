package cli

import (
	"os"
	"path"

	"github.com/feloy/kubernetes-api-reference/pkg/config"
	"github.com/feloy/kubernetes-api-reference/pkg/kubernetes"
	"github.com/spf13/cobra"
)

// ShowTOCCmd defines the `showtoc` subcommand
func ShowTOCCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "showtoc",
		Short:         "show the table of contents",
		Long:          "list the parts and chapter of the documentation",
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

			toc.ToMarkdown(os.Stdout)
			return nil
		},
	}
	cmd.Flags().StringP("config-dir", "c", "", "Directory conatining documentation configuration")
	cmd.MarkFlagRequired("config-fir")

	return cmd
}

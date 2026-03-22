/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cli

import (
	"log"
	"os"
	"strings"

	"github.com/keikoproj/aws-auth/pkg/mapper"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

var getArgs = &mapper.MapperArguments{
	OperationType: mapper.OperationGet,
	IsGlobal:      true,
}

// getCmd represents the base view command when run without subcommands
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get provides a detailed summary of the configmap",
	Long:  `get allows a user to output the aws-auth configmap entires in various formats`,
	Run: func(cmd *cobra.Command, args []string) {
		options := kubeOptions{
			AsUser:   upsertArgs.AsUser,
			AsGroups: upsertArgs.AsGroups,
		}

		k, err := getKubernetesClient(getArgs.KubeconfigPath, options)
		if err != nil {
			log.Fatal(err)
		}

		worker := mapper.New(k, true)

		d, err := worker.Get(getArgs)
		if err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewTable(os.Stdout,
			tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
				Borders: tw.Border{
					Left:   tw.Off,
					Right:  tw.Off,
					Top:    tw.Off,
					Bottom: tw.Off,
				},
				Settings: tw.Settings{
					Lines: tw.Lines{
						ShowTop:        tw.Off,
						ShowBottom:     tw.Off,
						ShowHeaderLine: tw.Off,
					},
					Separators: tw.Separators{
						BetweenColumns: tw.Off,
						BetweenRows:    tw.Off,
					},
				},
			})),
		)
		table.Configure(func(cfg *tablewriter.Config) {
			cfg.Header.Alignment.Global = tw.AlignLeft
			cfg.Header.Formatting.AutoFormat = tw.Off
			cfg.Row.Alignment.Global = tw.AlignLeft
			cfg.Header.Padding.Global = tw.Padding{Left: "\t", Right: "", Overwrite: true}
			cfg.Row.Padding.Global = tw.Padding{Left: "\t", Right: "", Overwrite: true}
		})
		table.Header([]string{"Type", "ARN", "Username", "Groups"})

		for _, row := range d.MapRoles {
			if err := table.Append([]string{"Role Mapping", row.RoleARN, row.Username, strings.Join(row.Groups, ", ")}); err != nil {
				log.Fatal(err)
			}
		}

		for _, row := range d.MapUsers {
			if err := table.Append([]string{"User Mapping", row.UserARN, row.Username, strings.Join(row.Groups, ", ")}); err != nil {
				log.Fatal(err)
			}
		}

		if err := table.Render(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVar(&getArgs.KubeconfigPath, "kubeconfig", "", "Path to kubeconfig")
	getCmd.Flags().StringVar(&getArgs.Format, "format", "table", "The format in which to display results (currently only 'table' supported)")
	getCmd.Flags().StringVar(&upsertArgs.AsUser, "as", "", "Username to impersonate for the operation")
	getCmd.Flags().StringSliceVar(&upsertArgs.AsGroups, "as-group", []string{}, "Group to impersonate for the operation, this flag can be repeated to specify multiple groups")
}

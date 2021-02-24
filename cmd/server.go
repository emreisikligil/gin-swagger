// Copyright © 2021 Emre Isikligil <emreisikligil@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/go-swagger/go-swagger/generator"
	"github.com/spf13/cobra"
)

type ServerOptions struct {
	Toggle         bool
	SkipOperations bool
}

var (
	// newCmd represents the new command
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Generates a new server",
		Long:  ``,
		Run: func(_ *cobra.Command, _ []string) {
			err := opts.EnsureDefaults()
			checkError(err)

			err = opts.CheckOpts()
			checkError(err)

			for _, asset := range AssetNames() {
				data, err := Asset(asset)
				checkError(err)

				err = generator.AddFile(asset, string(data))
				checkError(err)
			}

			err = generator.GenerateServer(opts.Name, nil, nil, &opts)
			checkError(err)

		},
	}
	opts    generator.GenOpts
	srvOpts ServerOptions
)

func init() {
	rootCmd.AddCommand(serverCmd)
	StringVarP := serverCmd.Flags().StringVarP
	BoolVarP := serverCmd.Flags().BoolVarP
	StringVarP(&opts.Name, "name", "n", "name", "name")
	StringVarP(&opts.APIPackage, "api-package", "p", "operations", "api-package name")
	StringVarP(&opts.ModelPackage, "model-package", "m", "models", "model-package name")
	StringVarP(&opts.ServerPackage, "server-package", "s", "restapi", "server-package name")
	StringVarP(&opts.ClientPackage, "client-package", "c", "client", "client-package name")
	StringVarP(&opts.Spec, "spec", "", "api/swagger.yml", "spec name")
	StringVarP(&opts.Target, "target", "", "./", "target name")
	BoolVarP(&opts.IncludeModel, "include-models", "", true, "Generates models if set")
	BoolVarP(&opts.IncludeHandler, "include-handler", "", true, "Generates operation handlers if set")
	BoolVarP(&opts.IncludeSupport, "include-support", "", true, "Generates support docs if set")
	BoolVarP(&opts.IncludeValidator, "include-validation", "", true, "Generates validators if set")
	BoolVarP(&opts.IncludeMain, "include-main", "", false, "Generates main function if set")
	BoolVarP(&opts.ExcludeSpec, "exclude-spec", "", false, "Excludes spec if set")
	BoolVarP(&srvOpts.SkipOperations, "skip-operations", "", false, "Skips operations if set")
	BoolVarP(&srvOpts.Toggle, "toggle", "t", false, "Help message for toggle")
	opts.Sections = defaultSectionOpts()
	generator.FuncMap["add"] = func(a, b int) int { return a + b }
}

func defaultSectionOpts() generator.SectionOpts {
	return generator.SectionOpts{
		Application: []generator.TemplateOpts{
			{
				Name:       "server",
				Source:     "templates/api.gotmpl",
				Target:     "{{ joinFilePath .Target .ServerPackage }}",
				FileName:   "server.go",
				SkipExists: false,
				SkipFormat: false,
			},
		},
		Operations: []generator.TemplateOpts{
			{
				Name:       "handler",
				Source:     "templates/operation.gotmpl",
				Target:     "{{ if gt (len .Tags) 0 }}{{ joinFilePath .Target .ServerPackage .APIPackage .Package  }}{{ else }}{{ joinFilePath .Target .ServerPackage .Package  }}{{ end }}",
				FileName:   "{{ (snakize (pascalize .Name)) }}.go",
				SkipExists: false,
				SkipFormat: false,
			},
		},
		Models: []generator.TemplateOpts{
			{
				Name:       "definition",
				Source:     "asset:model",
				Target:     "{{ joinFilePath .Target .ModelPackage }}",
				FileName:   "{{ (snakize (pascalize .Name)) }}.go",
				SkipExists: false,
				SkipFormat: false,
			},
		},
	}
}

func checkError(err error) {
	if err == nil {
		return
	}
	log.Fatalln("Error: %#v\n", err)
}

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/marcom4rtinez/infrahub-terraform-provider-generator/pkg/parser"
)

func main() {
	graphqlDirectory := flag.String("gql-dir", "gql", "Directory with GraphQL queries")
	providerDirectory := flag.String("provider-dir", "internal/provider", "Directory to write the generated Terraform Provider")

	flag.Parse()

	gqlDir := *graphqlDirectory

	var dataSources, resources []string

	err := filepath.Walk(gqlDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if filepath.Ext(path) == ".gql" {
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				dataSourceName, resourceName, err := parser.ReadAndGenerateDataSourcesAndResources(string(data), *providerDirectory)
				if err == nil {
					if dataSourceName != "" {
						dataSources = append(dataSources, dataSourceName)
					} else if resourceName != "" {
						resources = append(resources, resourceName)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	parser.ReadAndGenerateProvider(
		parser.TerraformComponents{
			DataSources: dataSources,
			Resources:   resources,
		}, *providerDirectory)
}

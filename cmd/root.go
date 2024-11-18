/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/aquasecurity/table"
	"github.com/gocolly/colly"
	"github.com/gookit/slog"
	"github.com/liamg/tml"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gojobs",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: findJobs,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gojobs.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Source struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Sources struct {
	Sources []Source `json:"sources"`
}


func findJobs(cmd *cobra.Command, args []string) {

	jsonFile, err := os.Open("sources.json")
	if err != nil {
		slog.Error(err)
		os.Exit(1)
	}
	defer jsonFile.Close()

	slog.Info("Successfully Opened sources.json")
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		slog.Error(err)
		os.Exit(1)
	}

	var sources Sources
	if err := json.Unmarshal(byteValue, &sources); err != nil {
		slog.Error(err)
		os.Exit(1)
	}

	sourcesMap := make(map[string]string)
	for _, source := range sources.Sources {
		sourcesMap[source.Name] = source.URL
	}

	slog.Info(sourcesMap)

	_  = colly.NewCollector()


	t := table.New(os.Stdout)
	t.SetPadding(1)
	t.SetAlignment(table.AlignCenter)
	t.SetAutoMerge(true)
	t.SetDividers(table.MarkdownDividers)


	t.SetHeaders("Title", tml.Sprintf("<green>Company</green>"), "Location", "Date", "Link")
	t.AddRow("Title", "Company", "Location", "Date", "Link")
	t.AddRow("Title", "Company", "Location", "Date", "[Link](https://google.com)")

	t.Render()
}

func createLink(linkText, url string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", url, linkText)
}
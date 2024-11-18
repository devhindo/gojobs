/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aquasecurity/table"
	"github.com/gocolly/colly"
	"github.com/gookit/slog"
	"github.com/liamg/tml"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	jobs []Job
)

type Job struct {
	Title    string
	Company  string
	Location string
	Date     string
	Description string
	Link     string
}

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

	slog.Print(jobs)

	c := colly.NewCollector()

	spinner, _ := pterm.DefaultSpinner.Start("Fetching jobs...")
	fetchGolangPrjects(c, sourcesMap["golangprojects"])
	spinner.RemoveWhenDone = true
	spinner.Success("Jobs fetched successfully!")

	t := table.New(os.Stdout)
	t.SetAlignment(table.AlignCenter)


	t.SetHeaders("Title", tml.Sprintf("<green>Company</green>"), "Location", "Date", "Description", "Link")
	i := 0
	for _, job := range jobs {
		i++
		if i == 3 {
			break
		}
		t.AddRow(job.Title, job.Company, job.Location, job.Location,job.Date, job.Description, job.Link)
	}

	t.Render()

	tableData := pterm.TableData{
		{"Title", "Company", "Location", "Date", "Link"},
	}
	for _, job := range jobs {
		tableData = append(tableData, []string{job.Title, job.Company, job.Date, job.Description, job.Link})
	}

	// Create a table with the defined data.

	// The table has a header and is boxed.
	// Finally, render the table to print it.
	pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(tableData).Render()

	fmt.Println(jobs[0])
}

// https://www.golangprojects.com/golang-go-job-gta-Staff-Back-End-Golang-Engineer-New-York-NY-NYC-ONRAMP.html





func fetchGolangPrjects(c *colly.Collector, source string) {
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
	/*
	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach("div.bg-csdpromobg1", func(_ int, el *colly.HTMLElement) {
			fmt.Println(el)

		})
	})
	*/

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		log.Println(link)
		if len(link) > 0 && link[0] == '/' && len(link) > 14 && link[:14] == "/golang-go-job" {
		//	log.Println(link)
			iTags := e.DOM.NextAllFiltered("i")
			iContents := []string{}
		//	log.Println(iContents)
			if iTags.Length() >= 2 {
				iTags = iTags.Slice(0, 2)
				iTags.Each(func(_ int, el *goquery.Selection) {
					iContents = append(iContents, el.Text())
				})
				log.Println(iContents)
			}
			
			if len(iContents) == 2 {
				bTag := e.DOM.Find("b").First()
				maintitle := bTag.Text()
				titleParts := strings.LastIndex(maintitle, "-")
				log.Println(titleParts)
				title := maintitle[:titleParts]
				company := maintitle[titleParts+1:]

				date := e.DOM.NextFiltered("i").Text()
				description := e.DOM.NextFiltered("i").Text()
				//log.Println("locationnnnnn")
				job := Job{
					Title:    title,
					Company: company,
					Location: iContents[1] + " ",
					Date:     date,
					Description: description,
					Link:    "https://www.golangprojects.com/" + link,
				}
				jobs = append(jobs, job)
			}
		}
	})
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL)
	})

	c.Visit(source)

	c.OnResponse(func(r *colly.Response) {
	//	fmt.Println("Visited", r.Request.URL)
	//	fmt.Println(string(r.Body))
	})

	slog.Println(jobs)
//	slog.Info(len(jobs))
}


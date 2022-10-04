package main

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	gaRep "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/option"
)

func main() {
	fmt.Println("starting connection...")
	config := oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	token := &oauth2.Token{
		AccessToken:  os.Getenv("ACCESS_TOKEN"),
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
	}

	ctx := context.Background()
	gaService, err := gaRep.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))

	if err !=nil {
		fmt.Println("error creating service:", err)
	}

	report, err := getReport(gaService)
	
	if (err !=nil) {
		fmt.Println("error creating report:", err)
	}
	fmt.Println("REPORT", report)

	printResponse(report)

}

func getReport(svc *gaRep.Service) (*gaRep.GetReportsResponse, error) {
	req := &gaRep.GetReportsRequest{
		// Our request contains only one request
		// So initialise the slice with one ga.ReportRequest object
		ReportRequests: []*gaRep.ReportRequest{
			// Create the ReportRequest object.
			{
				ViewId: "261593436",
				DateRanges: []*gaRep.DateRange{
					// Create the DateRange object.
					{StartDate: "7daysAgo", EndDate: "today"},
				},
				Metrics: []*gaRep.Metric{
					// Create the Metrics object.
					{Expression: "ga:users"},
					{Expression: "ga:sessions"},
				},
				Dimensions: []*gaRep.Dimension{
					{Name: "ga:country"},
				},
			},
		},
	}
	return svc.Reports.BatchGet(req).Do()
}

// printResponse parses and prints the Analytics Reporting API V4 response.
func printResponse(res *gaRep.GetReportsResponse) {
	for _, report := range res.Reports {
		header := report.ColumnHeader
		dimHdrs := header.Dimensions
		metricHdrs := header.MetricHeader.MetricHeaderEntries
		rows := report.Data.Rows

		if rows == nil {
			fmt.Println("No data found for given view.")
		}

		for _, row := range rows {
			dims := row.Dimensions
			metrics := row.Metrics

			for i := 0; i < len(dimHdrs) && i < len(dims); i++ {
				fmt.Printf("%s: %s", dimHdrs[i], dims[i])
				
			}

			for _, metric := range metrics {
				// We have only 1 date range in the example
				// So it'll always print "Date Range (0)"
				// log.Infof("Date Range (%d)", idx)
				for j := 0; j < len(metricHdrs) && j < len(metric.Values); j++ {
					 fmt.Printf(" %s: %s ", metricHdrs[j].Name, metric.Values[j])
				}
			}

			fmt.Println("")
		}
	}
	
}
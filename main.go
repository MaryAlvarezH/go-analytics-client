package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gaRep "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/option"
)

func main() {
	// Step 1. Setup credentials for oAuth authentication
	// ClientID ClientSecret -> https://console.cloud.google.com/apis/credentials/oauthclient/:PROJECT_ID
	// RedirectURL -> it's always "https://developers.google.com/oauthplayground" value
	// Endpoint -> AuthURL:  "https://provider.com/o/oauth2/auth", TokenURL: "https://provider.com/o/oauth2/token",
	// AccessToken & RefreshToken -> Use https://developers.google.com/oauthplayground to generate them using ClientID ClientSecret values in oAuth Configuration (readme*)
	// Expiry -> it's always with time.Now() to reuse the same token
	config := oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  os.Getenv("REDIRECT_URL"),
	}

	token := &oauth2.Token{
		AccessToken:  os.Getenv("ACCESS_TOKEN"),
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		Expiry:       time.Now(),
	}

	// Step 2. Create the service instance using analyticsreporting/v4 package
	ctx := context.Background()
	gaService, err := gaRep.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))

	if err != nil {
		fmt.Println("error creating service:", err)
	}

	// Step 3. Create report
	report, err := getReport(gaService)

	if err != nil {
		fmt.Println("error creating report:", err)
	}

	printResponse(report)

}

func getReport(svc *gaRep.Service) (*gaRep.GetReportsResponse, error) {
	req := &gaRep.GetReportsRequest{
		// Our request contains only one request
		// So initialise the slice with one ga.ReportRequest struct.
		ReportRequests: []*gaRep.ReportRequest{
			// Create the ReportRequest struct.
			{
				// ViewId -> Analytics viewId (readme*)
				ViewId: "261593436",
				DateRanges: []*gaRep.DateRange{
					// Create the DateRange struct.
					{StartDate: "2022-09-01", EndDate: "today"},
				},
				// Create the Metrics struct.
				Metrics: []*gaRep.Metric{

					{Expression: "ga:adClicks", Alias: "Clicks"},
					{Expression: "ga:sessions", Alias: "Sessions"},
					{Expression: "ga:users", Alias: "Users"},
					{Expression: "ga:transactions", Alias: "Transactions"},
					{Expression: "ga:transactionRevenue", Alias: "Revenue"},
				},
				// Create the Dimensions struct.
				Dimensions: []*gaRep.Dimension{

					{Name: "ga:campaign"},
					{Name: "ga:date"},
				},
				// Create the MetricFilterClauses struct.
				MetricFilterClauses: []*gaRep.MetricFilterClause{
					{
						Filters: []*gaRep.MetricFilter{
							{
								MetricName:      "ga:transactions",
								Operator:        "GREATER_THAN",
								ComparisonValue: "100",
							},
						},
					},
				},
			},
		},
	}

	// Use BatchGet method to return the analytics data
	return svc.Reports.BatchGet(req).Do()
}

// printResponse parses and prints the Analytics Reporting API V4 response.
func printResponse(res *gaRep.GetReportsResponse) {
	// Response body format -> https://developers.google.com/analytics/devguides/reporting/core/v4/advanced#response_column_header
	for _, report := range res.Reports {
		header := report.ColumnHeader
		dimHdrs := header.Dimensions                          // dimensions in headers
		metricHdrs := header.MetricHeader.MetricHeaderEntries // metrics in headers
		rows := report.Data.Rows                              // registers founds

		if rows == nil {
			fmt.Println("No data found for given view.")
		}

		// For each register it's metrics and dimensions are printed
		for _, row := range rows {
			dims := row.Dimensions
			metrics := row.Metrics

			// dimensions
			for i := 0; i < len(dimHdrs) && i < len(dims); i++ {
				fmt.Printf(" %s: %s ", dimHdrs[i], dims[i])

			}

			// metrics
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

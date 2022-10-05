# Connection and use of Analytics reporting API v4

## Connection

1. Create or use an existing oAuth credential from Google Cloud Platform > Specific project > API's & Services > Credentials > oAuth 2.0 Client ID's

   <img width="1442" alt="Screen Shot 2022-10-05 at 9 52 30" src="https://user-images.githubusercontent.com/38545126/194095690-4c89cc2a-1e6c-40a1-a295-992b1469e0e6.png">

   - ClientID and ClientSecret will be used in the config setup for client connection
     <img width="1505" alt="Screen Shot 2022-10-05 at 10 10 38" src="https://user-images.githubusercontent.com/38545126/194096301-89f3dfb4-2b92-4f9f-8b96-e6ab468ee683.png">

   - Add the following URI's to authorized redirect URI's https://developers.google.com/oauthplayground

2. Go to https://developers.google.com/oauthplayground

   - Click on oAuth 2.0 configuration check "Use your own oAuth credentials" and write your ClientID and ClientSecret.
     <img width="491" alt="Screen Shot 2022-10-04 at 17 40 18" src="https://user-images.githubusercontent.com/38545126/193944311-f5272b86-4372-4a33-bc4a-add4e528a326.png">

   - In Step 1 select Analytics Reporting API v4 and clicks on authorize API button
     <img width="1132" alt="Screen Shot 2022-10-05 at 10 19 40" src="https://user-images.githubusercontent.com/38545126/194098325-cf93653d-1c27-4c59-b1ac-e1502a693b6a.png">

   - Select your google account
   
     <img width="491" alt="Screen Shot 2022-10-05 at 10 19 55" src="https://user-images.githubusercontent.com/38545126/194098328-f86d9580-df5e-483f-b97b-cfe6c948d3c0.png">

   - Conceed the requested permissions
   
     <img width="484" alt="Screen Shot 2022-10-05 at 10 20 07" src="https://user-images.githubusercontent.com/38545126/194098335-11a9c1dc-9c43-4b51-b76b-1d94611656a8.png">

   - In Step 2 exchange your authorization code for tokens
     <img width="1131" alt="Screen Shot 2022-10-05 at 10 20 24" src="https://user-images.githubusercontent.com/38545126/194098340-67494fd7-bfe2-4b0c-ac3a-846812254e7b.png">

   - Save your access token and refresh token
   
     <img width="491" alt="Screen Shot 2022-10-05 at 10 20 42" src="https://user-images.githubusercontent.com/38545126/194098362-76b0b984-7368-49cb-a6c4-4d98542d2250.png">

3. Use https://pkg.go.dev/google.golang.org/api/analyticsreporting/v4 package to connect with Analytics reporting API

## Making requests

Request and response have an specific format, you can check them in the following [link](https://developers.google.com/analytics/devguides/reporting/core/v4/advanced)

For advances request you can use this [request composer](https://ga-dev-tools.web.app/request-composer/)

Request example using Golang client:

```
req := &gaRep.GetReportsRequest{
		ReportRequests: []*gaRep.ReportRequest{
			{
				ViewId: "00000000",
				DateRanges: []*gaRep.DateRange{
					{StartDate: "2022-09-01", EndDate: "today"},
				},
				Metrics: []*gaRep.Metric{
					{Expression: "ga:users", Alias: "Users"},
					{Expression: "ga:sessions", Alias: "Sessions"},
					{Expression: "ga:transactions", Alias: "Transactions"},
					{Expression: "ga:transactionRevenue", Alias: "Revenue"},
				},
				Dimensions: []*gaRep.Dimension{
					{Name: "ga:country"},
					{Name: "ga:segment"},
				},
				MetricFilterClauses: []*gaRep.MetricFilterClause{
					{
						Filters: []*gaRep.MetricFilter{
							{
								MetricName:      "ga:transactions",
								Operator:        "GREATER_THAN",
								ComparisonValue: "20",
							},
						},
					},
				},
				DimensionFilterClauses: []*gaRep.DimensionFilterClause{
					{
						Filters: []*gaRep.DimensionFilter{
							{
								DimensionName: "ga:country",
								Operator:      "EXACT",
								Expressions:   []string{"Mexico"},
							},
						},
					},
				},
				Segments: []*gaRep.Segment{
					{SegmentId: "gaid::-2"},
				},
			},
		},
	}
```

Response example using Golang client:

```
{
    "reports": [
        {
            "columnHeader": {
                "metricHeader": {
                    "metricHeaderEntries": [
                        {
                            "name": "ga:users",
                            "type": "INTEGER"
                        }
                    ]
                }
            },
            "data": {
                "isDataGolden": true,
                "maximums": [
                    {
                        "values": [
                            "98"
                        ]
                    }
                ],
                "minimums": [
                    {
                        "values": [
                            "98"
                        ]
                    }
                ],
                "rowCount": 1,
                "rows": [
                    {
                        "metrics": [
                            {
                                "values": [
                                    "98"
                                ]
                            }
                        ]
                    }
                ],
                "totals": [
                    {
                        "values": [
                            "98"
                        ]
                    }
                ]
            }
        }
    ]
}
```

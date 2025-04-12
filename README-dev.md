
## Google Sheet API
- https://canopas.com/golang-interacting-with-google-spreadsheets-b381f819c2eb
- https://developers.google.com/identity/protocols/oauth2/service-account
- https://developers.google.com/sheets/api/quickstart/go
- https://console.cloud.google.com/apis/dashboard?project=cvfs-jpx
Remember to share the spreadsheet with service account email (nhuongmh-jpx@cvfs-jpx.iam.gserviceaccount.com)


## Unit Test
go test -v -run TestMaziiFetcher_FetchMaziiKanji github.com/nhuongmh/cfvs.jpx/pkg/supporter/mazii


PostgreSQL in WSL:
user=postgres, pass=clover.fox (access through sudo -su)
user=clover, pass=foxie
psql -U clover -d ie

#: \dt (list table)
-l : list database
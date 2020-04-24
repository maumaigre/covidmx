# COVID-19 MX API / CSV Fetch

Covidmx is a server-side application that fetches official .gob.mx data served daily with every COVID-19 case and feeds them to a MySQL Database.
A CRON job (default: each hour) downloads the new CSV and diffs with the latest one fed to the database and patches new data.

An API is also provided with general stats (counts of confirmed, active, etc).

Built in Go (v1.14) with the following dependencies:
* Mux Router
* Squirrel Query Builder
* Crontab
* SQLx

## WIP
* Filter/Query cases
* Front-end app / graphs
* API Usage Readme

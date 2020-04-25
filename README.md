
# COVID-19 MX API / CSV Fetch
   ![API](https://i.imgur.com/zF9jvgA.jpg?style=center)
 

Covidmx is a server-side application that fetches official .gob.mx data served daily with every COVID-19 case and feeds them to a MySQL Database.
A CRON job (default: each hour) downloads the new CSV and diffs MD5 checksum with the latest one fed to the database and dumps the data to the DB.

An API is also provided with general stats (counts of confirmed, active, etc).
  
Built in Go (v1.14) with the following dependencies:

* Mux Router
* Squirrel Query Builder
* Crontab
* SQLx

## WIP
* Filter/Query cases (struct for possible filterable parameters)
* Block force fetch if current goroutine is in process and once every X minutes
* Front-end app / graphs
* API Usage readme
* Dynamic SQL query statements
* Use structs for SQL queries
* Figure out a better way to dump SQL (Current process is a WIP with a shell script)
# Gator
Gator is a basic command line RSS aggregator, use PostgreSQL as a database for data persistence.
## Requirements
The go programming language has to be installed first to be able to build the tool, go is available on most Linux distribution package managers, as an example on Debian/Ubuntu
``` bash
sudo apt update
sudo apt install go
```
PostgreSQL is also a dependency to be able to use this tool, can be installed using the commands:
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```
A configuration file named gatorconfig.json is required to exist first
``` json
{
  "db_url": "postgres://example"
}
```
The db_url has to follow the next structure
```
protocol://username:password@host:port/database?sslmode=disable
```
Once you have the necessary information you can create this file using this command, make sure to replace the HOST and port with yours:
``` bash
echo "{
  \"db_url\": \"postgres://postgres:postgres@HOST:PORT/gator?sslmode=disable\"
}" > $HOME/.gatorconfig.json
```
You can create the tables using the queries found in this repository.
Once you have completed these steps, you can install this tool using:
``` bash
go install github.com/afcaballero-1994/gator
```
## Usage
To register a new user, you can use the command ``` register``` like this:
``` bash
gator register <username>
```
To login or set a new user:
``` bash
gator login <username>
```
To add feeds to the database:
``` bash
gator addfeed <name> <link>
```
To scrape feeds added:
``` bash
gator agg <time>
```
Example:
``` bash
gator agg 1m
```
And to get the posts added from feeds:
``` bash
gator browsw <limit>
```
The default is 2, and this limit is used in case an invalid limit is used

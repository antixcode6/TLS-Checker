# TLS-Checker
Utilizing a copy of rapidloop/certchk to begin interally check for and alert on valid SSL certs. 

# Usage
The script can take a single or multiple urls as params, or can take a file as a param. To run with the example url list provided you can do it like below

`go run tlscheck.go -f examples.txt`

<p>Additionally you can build the Go file and run it</p>
./tlscheck -f examples.txt 
<p></p>
This program will write any errors to a log file. On my list of things to do is to create an alert module so when errors are detected the log file can be emailed out to whomever is defined in the script (or from an email group) as needing to be alerted of an invalid, or soon to expire cert. 

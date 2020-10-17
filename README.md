# TLS-Checker
Utilizing a copy of rapidloop/certchk to begin interally check for and alert on valid SSL certs. 
<p></p>

# What it does
This program checks for valid SSL certs on urls you put into the input file, or if you want to run against a single (or list) of urls from the command line you can. It will write any errors to a log file and can be configured to send an email with the log file.

# How to use it
The program can take a single or multiple urls as params, or can take a file as a param. To run with the example url list provided you can do it like below

`go run tlscheck.go -f examples.txt`

<p>You can build the Go file and run it as there is a significant performance boost from it</p>

`go build tlscheck.go -> ./tlscheck -f examples.txt`

# Benchmarks
As I am learning Go while working on this project I wanted to test the program against a list of 2000 urls. 

Here are some stats I collected while hacking. 

Benchmark collected very scientifically by running

`time ./tlscheck -f files.txt`


```
Compiled and ran against 2000 host names with no goroutines
	Executed in  720.21 secs   fish           external
	   usr time   11.71 secs  543.00 micros   11.71 secs
	   sys time    4.03 secs  272.00 micros    4.03 secs

Not compiled and ran against 2000 host names with goroutines
	Executed in  640.90 secs   fish           external
	   usr time   10.79 secs  161.00 micros   10.79 secs
	   sys time    3.45 secs   57.00 micros    3.45 secs

Compiled and ran against 2000 with goroutines
	Executed in  495.48 secs   fish           external
	   usr time   11.06 secs  518.00 micros   11.06 secs
	   sys time    3.53 secs  257.00 micros    3.53 secs
```

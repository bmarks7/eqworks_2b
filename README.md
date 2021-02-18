# eqworks_2b
Answer to question 2b of product work sample from eq works

Once you have this repository on your computer, using your terminal/command line, enter the directory with the file named "main.go", and then type "go run main.go" and the server will begin. To implement a global rate limit on the "/stats" endpoint, I implemented a fixed window algorithm, where for every 10 second interval, a maximum of 5 calls to the endpoint are permitted.

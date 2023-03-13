### Description
This simple Consumer Loan Application is a web application that allows customer to calculate the monthly payment for a consumer loan. The application is written in GO and only uses built-in functions.

### Usage: how to run
```go run .``` to start the server at port 8080.
You can now visit the URL: http://localhost:8080/ or directly click the link in terminal.

### Implementation
The program calculates the given values and returns the monthly payment based on an interest rate of 5%. It also checks if the borrower is eligible for the loan. The borrower is eligible if he/she is not blacklisted and has not made more than 5 applications in the last 24 hours.
<br><br>
The file ``blacklist.txt`` contains a list of blacklisted customers.<br>
The file ``applications.txt`` logs all applications made by customers with timestamps.
<br><br>
Docker: The application can be run in a docker container. To build the image, run the following command in the project directory:
```docker build -t consumer-loan-app .```
followed by:<br>
```docker run -p 8080:8080 consumer-loan-app```

### Testing
There are also testfiles for the neccessary functions. To run the tests, run the following command in the project directory:
```go test```

### Final thoughts
This is my first full web application written in GO. I learned a lot about back-end communication and web development in general.

### Authors
- Alexander Embacher (https://github.com/4stroPhysik3r)
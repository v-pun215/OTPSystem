# OTP System
A simple OTP system that generates 6-digit passcode that expires in 5 minutes for an email.

## Parameters
 - /generate - email field required, generated 6 digit passcode for given email
 - /validate - email, code field required, validates code given for given email (also checks if done within time limit)

## Usage
 - Run server using ```go run main.go```.
 - Interface with it using ```curl "https://localhost:8080/...``` (use parameters)
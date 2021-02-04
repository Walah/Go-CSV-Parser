# CSV Parser in Go

I was working on a project where we had multiple services with in one project (e.g. 
Web, SafeKeeping, IOS, Call Notice, etc.). **WEB** was the front-end that will 
let the user upload the CSV file, and prepare it. After that, other services can 
collect the data using a POST request.

The whole CSV file parsing is devided into 2 parts...
  - MultipartWriter() 
  - MultipartReader() 

## MultipartWriter 

It converts the contents of the passed/received file into a `multipart.Writer` 
form file so the data can be passed to next service.

It returns the correct HTTP Content-Type, the form content as bytes, 
and a possible error.

(**WEB** service will collect the CSV file from the user, prepare it using 
`MultipartWriter()` and send it off to other/next service.)

## MultipartReader 

It reads the records from the passed/received CSV-encoded file and 
returns the parsed data (in a 2D String slice format) to the next service.

It also gives an option to include or exclude file header. 

(Other services will collect the file/data from the POST request and 
parse it using the `MultipartReader()`)

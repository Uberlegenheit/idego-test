# idego-test
Test project for Idego Group

# Explanation
There are 2 versions of this task. You can see them in commits history.

## V1
Use `go run main.go <city name> <duration>` in root folder of project.
Example: `go run main.go Oregon 5m30s`. Minimal duration is

## V2
Use `go run main.go <city names separated by comma> <duration>` in root folder of project.
Example: `go run main.go Oregon,Miami,Dallas 5m30s`. Minimal duration is 10s.
In this version I used channels.

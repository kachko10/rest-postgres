package main

import "os"

func main() {
	a := App{}
	os.Setenv("APP_DB_USERNAME", "admin")
	os.Setenv("APP_DB_PASSWORD", "manager")
	os.Setenv("APP_DB_PASSWORD", "admin")

	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	a.Run(":8081")
}

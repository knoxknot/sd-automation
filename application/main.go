package main

import (
	"flag"
	"gitlab.com/knoxknot/csproject/application/controllers"
)


// program entry point
func main() {
	flag.Parse()

	//Start the API Server
	controllers.ServeAPI()
}
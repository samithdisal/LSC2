package main

//region Boiler plate
import (
	"flag"
	"log"

	"malhora.info/lew/crw"
)

func main() {
	log.Println("Starting LEW")

	url := flag.String("u", "", "URL to fetch")
	isServer := flag.Bool("s", false, "Should start as a server")
	serverHost := flag.String("h", "0.0.0.0:8080", "Start a server on given location")
	isVerbose := flag.Bool("v", false, "Verbose")
	isAuthor := flag.Bool("a", false, "Is author page")

	flag.Parse()

	// validate
	if !*isServer {
		if *url == "" {
			log.Fatal("url cannot be empty")
		}
	}
	//

	if *isServer {
		serv(*serverHost)
	} else {

		if *isVerbose {
			log.Println("Creating download requests for urls")
		}

		if *isAuthor {
			crw.GetAuthor(*url, "")
		} else {
			crw.GetPub(*url, "")
		}
	}
}

//endregion

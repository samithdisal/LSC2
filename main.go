package main

//region Boiler plate
import (
	"flag"
	"log"

	"malhora.info/lew/crw"
)

func main() {
	log.Println("Starting LEW")

	isServer := flag.Bool("s", false, "Should start as a server")
	serverHost := flag.String("h", ":8086", "Start a server on given location")
	isVerbose := flag.Bool("v", false, "Verbose")
	isAuthor := flag.Bool("a", false, "Is author page")

	flag.Parse()

	urls := flag.Args()

	if *isServer {
		serv(*serverHost)
	} else {

		if *isVerbose {
			log.Println("Creating download requests for urls")
		}

		if *isAuthor {
			for _, e := range urls {
				crw.GetAuthor(e, "")
			}
		} else {
			for _, e := range urls {
				crw.GetPub(e, "")
			}
		}
	}
}

//endregion

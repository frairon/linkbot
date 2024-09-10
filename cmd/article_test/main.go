package main

import (
	"log"
	"strings"

	"github.com/astappiev/microdata"
)

func main() {
	// Pass a URL to the `ParseURL` function.
	data, err := microdata.ParseURL("https://www.sueddeutsche.de/meinung/haushalt-ampel-lindner-scholz-kommentar-lux.3iCfPPE4tC9Au1TydFjDBE?reduced=true")
	if err != nil {
		log.Fatalf("error parsing: %v", err)
	}
	// https://schema.org/docs/full.html
	for _, item := range data.Items {
		if strings.Contains(item.Types[0], "Article") {
			headline, ok := item.Properties["headline"]
			if ok {
				log.Printf("%s", headline[0])
				break
			}

			for key, value := range item.Properties {
				log.Printf("%s: %#v", key, value)
			}
			return
		}
		log.Printf("%#v", item.Types)
	}

}

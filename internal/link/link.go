package link

import (
	"fmt"
	"strings"

	"github.com/astappiev/microdata"
)

type Link struct {
	Url      string
	Headline string
}

func parseLink(url string) (*Link, error) {

	var l Link

	l.Url = url

	data, err := microdata.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %w", err)
	}
	// https://schema.org/docs/full.html
	for _, item := range data.Items {
		if strings.Contains(item.Types[0], "Article") {
			headline, ok := item.Properties["headline"]
			if ok {
				// log.Printf("%s", headline[0])
				l.Headline = headline[0].(string)
				break
			}

			// for key, value := range item.Properties {
			// 	 log.Printf("%s: %#v", key, value)
			// }
			break
		}
		// log.Printf("%#v", item.Types)
	}
	return &l, nil
}

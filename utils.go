package requests

import (
	"fmt"
	"net/http"
	nUrl "net/url"
	"regexp"
	"strings"
)

const Version = "v1.1.15"

const DefaultTimeout = 5 // time.Second

var linkRegexCompiled = regexp.MustCompile(`(?:"|')(((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|((?:/|\.\./|\./)[^"'><,;|*()(%%$^/\\\[\]][^"'><,;|()]{1,})|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|/][^"|']{0,}|))|([a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:\?[^"|']{0,}|)))(?:"|')`)

func defaultHeaders() *http.Header {
	headers := &http.Header{}
	headers.Add("User-Agent", fmt.Sprintf("sari3l/requests %s", Version))
	return headers
}

func processLinks(url *nUrl.URL, links *[]string) *[]string {
	for index, link := range *links {
		link = strings.Trim(link, "\"")
		link = strings.Trim(link, "'")
		if len(link) >= 2 && link[0:2] == "//" {
			(*links)[index] = url.Scheme + ":" + link
		} else if len(link) >= 4 && link[0:4] == "http" {
			continue
		} else if len(link) >= 2 && link[0:2] != "//" {
			if link[0:1] == "/" {
				(*links)[index] = url.Scheme + "://" + url.Host + link
			} else if link[0:1] == "." {
				if link[0:2] == ".." {
					(*links)[index] = url.Scheme + "://" + url.Host + link[2:]
				} else {
					(*links)[index] = url.Scheme + "://" + url.Host + link[1:]
				}
			} else {
				(*links)[index] = url.Scheme + "://" + url.Host + "/" + link
			}
		} else {
			continue
		}
	}
	return links
}

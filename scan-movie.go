package gomovie

import (
	"strings"
	"net/http"
	"io/ioutil"
	"net/url"
	"errors"
)

func GetCode(title string) (string, error) {
	sites, err := googleSearch("Watch " + title + " Putlocker")
	if err != nil {
		return "", err
	}

	resp, err := http.Get(sites[0])
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile("site.html", body, 0644)
	if err != nil {
		return "", err
	}

	code, err := StringBetween(string(body), `<div class="video">
<script type="text/javascript">document.write(doit('`, `'));`)
	if err != nil {
		return "", err
	}

	return Decrypt(code), nil
}

// googleSearch searches a query to google.com and returns all
// the website url's on the first page in a slice of string.
func googleSearch(query string) (results []string, err error) {
	results = make([]string, 0, 10)
	query = strings.Replace(query, " ", "%20", -1)
	resp, err := http.Get("https://www.google.com.au/search?q="+query)
	if err != nil {
		return results, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return results, err
	}

	sites := strings.Split(string(body), `class="g"`)
	for _, site := range sites {
		if strings.Index(site, `<a href="`) != -1 {
			site = site[strings.Index(site, `<a href="`) + len(`<a href="`):]
			if site[:len(`/url?q=`)] == `/url?q=` {
				if strings.Index(site, `">`) != -1 {
					site = site[len(`/url?q=`):strings.Index(site, `">`)]
					if strings.Index(site, `&`) != -1 {
						site = site[:strings.Index(site, `&`)]
						site, err := url.QueryUnescape(site)
						if err != nil {
							return results, err
						}
						results = append(results, site)
					}
				}
			}
		}
	}

	return results, nil
}

// isPutlockerOnline checks weather a putlocker url's response
// and returns a bool representing weather it is accessable.
// An error is returned if a process failed during the process.
func isPutlockerOnline(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	str, err := StringBetween(string(body), `<title>`, `</title>`)
	if err != nil {
		return false, err
	}

	if strings.Index(strings.ToLower(str), strings.ToLower(`Website is offline`)) != -1 {
		return false, nil
	}

	return true, nil
}

// urlIsPutlock returns a bool showing weather
// the provided URL is located at putlocker;
// no matter the final url extension.
func urlIsPutlocker(url string) bool {
	splitted := strings.Split(url, ".")

	for i, section := range splitted {
		if strings.Index(section, "/") != -1 && i > 0 {
			if strings.ToLower(splitted[i-1]) == "putlocker" {
				return true
			}
		}
	}

	return false
}

// stringBetween returns a substring located between the first occurrence of
// both the provided start and end strings. An error will be returned if
// str does not include both start and end as a substring.
func StringBetween(str, start, end string) (string, error) {
	if strings.Index(str, start) == -1 || strings.Index(str, end) == -1 {
		return "", errors.New("String does not include start/end as substring.")
	}
	str = str[len(start):]
	return str[:strings.Index(str, end)], nil
}
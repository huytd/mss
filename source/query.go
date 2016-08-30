package source

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func GetHTTP(inputUrl string) *http.Response {
	ch := make(chan *http.Response)
	go func() {
		resp, err := http.Get(inputUrl)
		if err != nil {
			ch <- nil
		}
		ch <- resp
	}()
	for {
		select {
		case response := <-ch:
			return response
		}
	}
	return nil
}

func GetContent(inputUrl string) string {
	response := GetHTTP(inputUrl)
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "ERR"
	}
	defer response.Body.Close()
	html := string(b[:])
	return html
}

func ParseRegEx(input string, match string) string {
	re := regexp.MustCompile(match)
	matches := re.FindStringSubmatch(input)
	ret, err := url.QueryUnescape(matches[1])
	if err != nil {
		return "ERR"
	}
	return ret
}

func GetURL(inputUrl string) string {
	ch := make(chan string)
	if strings.Contains(inputUrl, "chiasenhac") {
		go func() {
			returnUrl := ChiaSeNhacURL(inputUrl)
			ch <- returnUrl
		}()
	} else if strings.Contains(inputUrl, "nhaccuatui") {
		go func() {
			returnUrl := NhacCuaTuiURL(inputUrl)
			ch <- returnUrl
		}()
	}
	for {
		select {
		case result := <-ch:
			return result
		}
	}
	return ""
}

func ChiaSeNhacURL(inputUrl string) string {
	html := GetContent(inputUrl)
	return ParseRegEx(html, `decodeURIComponent\(\"(.*\.m4a)\"\)`)
}

func NhacCuaTuiURL(inputUrl string) string {
	html := GetContent(inputUrl)
	xmlUrl := ParseRegEx(html, `\"(http:\/\/www\.nhaccuatui\.com\/flash\/xml.*)\"`)
	if xmlUrl != "ERR" {
		xml := GetContent(xmlUrl)
		return ParseRegEx(xml, `CDATA\[(http:\/\/.*)\]\]`)
	}
	return "ERR"
}

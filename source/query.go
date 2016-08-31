package source

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type (
	Cache map[string]string
)

var (
	CachedURL = make(Cache)
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
	if response != nil {
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "ERR"
		}
		defer response.Body.Close()
		html := string(b[:])
		return html
	} else {
		return "ERR"
	}
}

func ParseRegEx(input string, match string) string {
	re := regexp.MustCompile(match)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return "ERR"
	}
	ret, err := url.QueryUnescape(matches[1])
	if err != nil {
		return "ERR"
	}
	return ret
}

func GetURL(inputUrl string) string {
	ch := make(chan string)

	cachedUrl, isCached := CachedURL[inputUrl]
	if isCached {
		log.Print("Found cached URL: ", cachedUrl)
		return cachedUrl
	} else {
		log.Print("No cached URL found. Querying...")

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
		} else if strings.Contains(inputUrl, "mp3.zing") {
			go func() {
				returnUrl := ZingMp3URL(inputUrl)
				ch <- returnUrl
			}()
		}
	}
	for {
		select {
		case result := <-ch:
			CachedURL[inputUrl] = result
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

func ZingMp3URL(inputUrl string) string {
	songId := ParseRegEx(inputUrl, `http\:\/\/mp3\.zing\.vn\/bai\-hat\/.*\/(.*)\.html`)
	html := GetContent("http://api.mp3.zing.vn/api/mobile/song/getsonginfo?requestdata={\"id\":\"" + songId + "\"}")
	songUrl := ParseRegEx(html, `\"320\"\:\"(http\:\\\/\\\/api\.mp3\.zing\.vn\\\/api\\\/mobile\\\/source\\\/song\\\/\S[^\"]*)\"`)
	return strings.Replace(songUrl, `\`, "", -1)
}

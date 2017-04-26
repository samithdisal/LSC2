package main

//region Boiler plate
import (
	"flag"
	"log"
	"strings"

	"fmt"
	"regexp"

	"os"

	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	log.Println("Starting LEW")

	isVerbose := flag.Bool("v", false, "Verbose")
	isAuthor := flag.Bool("a", false, "Is author page")

	flag.Parse()

	urls := flag.Args()

	if *isVerbose {
		log.Println("Creating download requests for urls")
	}

	if *isAuthor {
		for _, e := range urls {
			getAuthor(e)
		}
	} else {

		for _, e := range urls {
			getPub(e)
		}
	}

}

//endregion

var pauseDuration time.Duration = 12 * time.Second

func getContent(url string) string {
	log.Println("Fetcing ", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println("Cannot fetch ", url, err)
		return ""
	}

	content := doc.Find("div.b-story-body-x.x-r15").First()
	title := doc.Find("div.b-story-header h1").First().Text()
	author := doc.Find("span.b-story-user-y.x-r22 a").First().Text()
	fileName := fmt.Sprintf("%s_%s.txt", Marshal(author, true), Marshal(title, true))
	out, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Print("Failed to create/open file ", fileName, err)
		return ""
	}
	stats, err := out.Stat()
	if err != nil {
		log.Print("Failed to read stats", err)
		return ""
	}
	if stats.Size() < 10 {
		cdate := time.Now().Format("2006 Jan 2")
		out.WriteString(fmt.Sprintf("%s by %s\n(%s)\nFetched on (%s)\n\n\n", title, author, url, cdate))
	}
	out.WriteString(content.Text())
	out.WriteString("\n\n---------------------------------\n\n")
	out.Close()
	log.Println("Done ", url)

	next := doc.Find("a.b-pager-next")
	if next.Length() > 0 {
		nxt_ret_val, exits := next.First().Attr("href")
		if exits {
			return nxt_ret_val
		}
	}
	writeEnd(fileName)
	return ""
}

func writeEnd(fileName string) {
	out, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Print("Failed to create/open file ", fileName, err)
	}
	out.WriteString("\n\n------------- THE END --------------------\n\n")
	out.Close()
}

func getPub(url string) {
	next := url
	for next != "" {
		next = getContent(next)
		time.Sleep(pauseDuration)
	}
}

func getAuthor(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println("Cannot fetch author page", err)
		return
	}
	pubs := doc.Find("tr.sl td a.bb").Map(getURLs)
	tot := len(pubs)
	for i, pub := range pubs {
		log.Printf("\n Fetching %d of %d\n", i, tot)
		getPub(pub)
	}
	log.Println("Done fetching author")
}

func getURLs(index int, se *goquery.Selection) string {
	return se.AttrOr("href", "")
}

//region Utilities

// Replacement structure
type replacement struct {
	re *regexp.Regexp
	ch string
}

// Build regexps and replacements
var (
	rExps = []replacement{
		{re: regexp.MustCompile(`[\xC0-\xC6]`), ch: "A"},
		{re: regexp.MustCompile(`[\xE0-\xE6]`), ch: "a"},
		{re: regexp.MustCompile(`[\xC8-\xCB]`), ch: "E"},
		{re: regexp.MustCompile(`[\xE8-\xEB]`), ch: "e"},
		{re: regexp.MustCompile(`[\xCC-\xCF]`), ch: "I"},
		{re: regexp.MustCompile(`[\xEC-\xEF]`), ch: "i"},
		{re: regexp.MustCompile(`[\xD2-\xD6]`), ch: "O"},
		{re: regexp.MustCompile(`[\xF2-\xF6]`), ch: "o"},
		{re: regexp.MustCompile(`[\xD9-\xDC]`), ch: "U"},
		{re: regexp.MustCompile(`[\xF9-\xFC]`), ch: "u"},
		{re: regexp.MustCompile(`[\xC7-\xE7]`), ch: "c"},
		{re: regexp.MustCompile(`[\xD1]`), ch: "N"},
		{re: regexp.MustCompile(`[\xF1]`), ch: "n"},
	}
	spacereg       = regexp.MustCompile(`\s+`)
	noncharreg     = regexp.MustCompile(`[^A-Za-z0-9-]`)
	minusrepeatreg = regexp.MustCompile(`\-{2,}`)
)

// Marshal function returns slugifies string "s"
func Marshal(s string, lower ...bool) string {
	for _, r := range rExps {
		s = r.re.ReplaceAllString(s, r.ch)
	}

	if len(lower) > 0 && lower[0] {
		s = strings.ToLower(s)
	}
	s = spacereg.ReplaceAllString(s, "-")
	s = noncharreg.ReplaceAllString(s, "")
	s = minusrepeatreg.ReplaceAllString(s, "-")

	return s
}

//endregion

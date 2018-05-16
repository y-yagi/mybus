package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func msg(err error, errStream io.Writer) int {
	if err != nil {
		fmt.Fprintf(errStream, "%v\n", err)
		return 1
	}
	return 0
}

func trimSpaces(text string) string {
	re := regexp.MustCompile(`[[:space:]]*`)
	return re.ReplaceAllString(text, "")
}

func run(args []string, outStream, errStream io.Writer) int {
	urls := map[string]string{
		"西口 ｔｏｃｏ川岸循環":    "http://transfer.navitime.biz/5931bus/pc/diagram/BusDiagram?orvCode=00020553&course=0001000836&stopNo=1",
		"西口 西川口駅西口ゆき":    "http://transfer.navitime.biz/5931bus/pc/diagram/BusDiagram?orvCode=00020553&course=0001000106&stopNo=1",
		"東口 蕨駅西口ゆき":      "http://transfer.navitime.biz/5931bus/pc/diagram/BusDiagram?orvCode=00020627&course=0001000349&stopNo=1",
		"中町二丁目 ｔｏｃｏ川岸循環": "http://transfer.navitime.biz/5931bus/pc/diagram/BusDiagram?orvCode=00021745&course=0001000836&stopNo=14",
		"中町二丁目 戸田公園駅西口":  "http://transfer.navitime.biz/5931bus/pc/diagram/BusDiagram?orvCode=00020482&course=0001000367&stopNo=6",
		"中町二丁目 西川口駅西口":   "http://transfer.navitime.biz/5931bus/pc/diagram/BusDiagram?orvCode=00020482&course=0001000242&stopNo=12",
		"中町二丁目 川口駅西口":    "http://transfer.navitime.biz/5931bus/pc/diagram/BusDiagram?orvCode=00020482&course=0001000397&stopNo=16",
	}

	for name, url := range urls {
		res, err := http.Get(url)
		if err != nil {
			fmt.Printf("Page get error: %v.\n", err)
			os.Exit(1)
		}

		defer res.Body.Close()
		if res.StatusCode != 200 {
			fmt.Printf("Status code error: %d %s", res.StatusCode, res.Status)
			os.Exit(1)
		}

		fmt.Fprintf(outStream, "%s: \n", name)
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return msg(err, errStream)
		}

		targetClasses := []string{".wkd", ".std", ".snd"}
		doc.Find(".diagram-table tr.l2").Each(func(_ int, s *goquery.Selection) {
			for _, class := range targetClasses {
				s.Find(class).Each(func(_ int, td *goquery.Selection) {
					td.Find(".diagram-item").Each(func(_ int, item *goquery.Selection) {
						text := trimSpaces(item.Find(".mm").Text())
						fmt.Fprintf(outStream, "%s  ", text)
					})
					fmt.Fprintf(outStream, "\t")
				})
			}
			fmt.Fprintf(outStream, "\n")
		})

		time.Sleep(1 * time.Second)
	}

	return 0
}

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}

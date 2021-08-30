package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)


func goGet() {
	var headings, row []string
	var rows [][]string

	data := `<html><body>
	<table>
		<tr><th>Heading 1</th><th>Heading two</th></tr>
		<tr><td>Data 11</td><td>Data 12</td></tr>
		<tr><td>Data 21</td><td>Data 22</td></tr>
		<tr><td>Data 31</td><td>Data 32</td></tr>
		<tr><td>Data 41</td><td>Data 42</td></tr>
	</table>
	<p>Stuff in here</p>
	<table>
		<tr><th>Heading 21</th><th>Heading 2two</th></tr>
		<tr><td>Data 211</td><td>Data 212</td></tr>
		<tr><td>Data 221</td><td>Data 222</td></tr>
		<tr><td>Data 231</td><td><span></span><span><a href="">Data 232</a></span></td></tr>
		<tr><td>Data 241</td><td>Data 242</td></tr>
	</table>
	</body>
	</html>
	`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		fmt.Println("No url found")
		log.Fatal(err)
	}

	// Find each table
	doc.Find("table").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			rowhtml.Find("th").Each(func(indexth int, tableheading *goquery.Selection) {
				headings = append(headings, tableheading.Text())
			})
			rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
				row = append(row, tablecell.Text())
			})
			rows = append(rows, row)
			row = nil
		})
	})
	fmt.Println("####### headings = ", len(headings), headings)
	fmt.Println("####### rows = ", len(rows), rows)
}

func main() {
	goGet()
}
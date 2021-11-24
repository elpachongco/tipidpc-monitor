// Summary
// Keeps watch of the query and prints new items when a new listing
// has been found.
//
// Usage
// To use, put a link in the queryUrl. The URL should be any page that will
// display a search query on tipidpc.
// the program runs every 5 seconds (default) but can be configured by replacing
// the 5 in the variable `interval`
//
// Bugs
// Some pages can't be fully read when the user's email is hidden.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	// queryUrl = `https://tipidpc.com/itemsearch.php?sec=s&namekeys=%22thinkpad%22`
	queryUrl = `https://tipidpc.com/catalog.php?cat=0&sec=s`
	interval = 5 * time.Second
)

type Item struct {
	Price string
	Name  string
	// Date  string
	// Url   string
}

func main() {

	items := []Item{}
	for oldItems := []Item{}; ; oldItems = items {
		log.Println("Searching...")
		items = Getlistings(queryUrl)

		// Compare old and new slice of listings.
		// Return a slice of new listings.
		newItems := CompareItems(oldItems, items)

		if len(oldItems) < 1 {
			Notify(newItems, "print")
		}

		time.Sleep(interval)
	}
}

func Getlistings(url string) []Item {
	// Send a Get request to server
	// Parse html. Put listings into a slice of
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal("Get Request Error")
	}
	defer resp.Body.Close()

	items, err := ParseListings(resp.Body)
	if err != nil {
		log.Fatal("Parse Error")
	}
	return items
}

// ParseListings scans the html data and puts listings into []Item
func ParseListings(d io.ReadCloser) ([]Item, error) {

	doc, err := html.Parse(d)
	if err != nil {
		return []Item{}, err
	}

	i := LocateUl(doc)

	return i, nil
}

// LocateUl scans the html and finds the <ul> that contains the search
// results table.
func LocateUl(t *html.Node) []Item {

	var items []Item
	if t.Data == "ul" {
		for _, v := range t.Attr {
			if v.Key == "id" && v.Val == "item-search-results" {
				items = ItemizeList(t)
				return items
			}
		}
	}

	for t := t.FirstChild; t != nil; t = t.NextSibling {
		r := LocateUl(t)
		for _, v := range r {
			items = append(items, v)
		}
	}
	return items

}

// ItemizeList scans each <li> and transforms it into an instance of
// type Item
func ItemizeList(t *html.Node) []Item {
	var items []Item
	var s []string
	// for each <li>
	for t := t.FirstChild; t != nil; t = t.NextSibling {
		recurse(t, &s)
	}

	// The item name appears in s[9], s[29], s[49], s[69], ...
	// The Price appears s[9 + 4], s[29 + 4], s[49 + 4], s[69 + 4]
	for k, v := range s {
		j := (k - 9) % 20 // Produces 0 every 9, 29, 49, 69
		if j == 0 {
			i := Item{Name: v, Price: s[k+4]}
			items = append(items, i)
		}
	}
	return items
}

func recurse(t *html.Node, s *[]string) {

	for t := t.FirstChild; t != nil; t = t.NextSibling {
		*s = append(*s, t.Data)
		recurse(t, s)
	}
}

// CompareItems returns []Item of Items new in curr based on prev.
// It works by looping through prev and curr, anything not found in prev is
// considered new.
// If loop encounters >= maxSeq matches, CompareItems returns to prevent looking up
// Older listings that were moved to a different page.
// Assumes that new items will be added in front of curr
func CompareItems(prev, curr []Item) []Item {

	matchIndex := []int{}
	// For each Item in Prev, loop through curr
	for _, v := range prev {
		for i, j := range curr {
			// Save the index if a match is found.
			if v == j {
				matchIndex = append(matchIndex, i)
			}
		}
	}

	if len(matchIndex) < 1 {
		// all items are new
		return curr
	}

	// Scan matches to see if they're in sequence
	for i := 0; i < len(matchIndex)-1; i++ {
		if matchIndex[i+1] != matchIndex[i]+1 {
			return []Item{}
		}
	}

	return curr[0:matchIndex[0]]
}

// Notify looks at all the items and sends a notification based on the mode
func Notify(i []Item, method string) {
	switch method {
	case "dmenu":
		var s string
		for k, v := range i {
			fmtStr := fmt.Sprintf("%d %.45s\t \t%.6s\n", k, v.Name, v.Price)
			s += fmtStr
			fmt.Print(fmtStr)
		}

		cmd := exec.Command("dmenu", "-l10")
		// Pipe the strings to the command
		cmd.Stdin = strings.NewReader(s)
		cmd.Run()

	default:
		for k, v := range i {
			fmt.Printf("%d %.45s\t \t%.6s\n", k, v.Name, v.Price)
		}
	}
}

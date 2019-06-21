package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"sync"
)

var wg sync.WaitGroup


func buildURL(word string) string {
	return "https://www.dictionary.com/browse/" + string(word)
}

func get(url string) []string {
	c := colly.NewCollector()
	c.IgnoreRobotsTxt = true
	var ret []string
	c.OnHTML("a.css-cilpq1.e15p0a5t2", func(e *colly.HTMLElement) {
		ret = append(ret, string(e.Text))
	})
	c.Visit(url)
	c.Wait()

	return ret
}

func threading(c chan []string, word string) {
	defer wg.Done()
	var words []string
	for _, w := range get(buildURL(word)) {
		words = append(words, w)
	}
	c <- words
}

func main() {
	fmt.Println("START")
	word := "jump"
	maxDepth := 2

	//bfs
	var q map[string]int
	nq := map[string]int {
		word: 0,
	}

	vis := make(map[string]bool)

	for i := 1; i <= maxDepth; i++ {
		fmt.Println(i)
		q, nq = nq, make(map[string]int)
		queue := make(chan []string, 5000)
		for word := range q {
			if _, ok := vis[word]; !ok {
				wg.Add(1)
				vis[word] = true
				go threading(queue, word)
			}
		}
		wg.Wait()
		close(queue)

		for v := range queue {
			fmt.Println(v)
			for _, w := range v {
				nq[w] = i
			}
		}
	}

	fmt.Println("END")
}

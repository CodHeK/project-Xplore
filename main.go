package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"sort"
	"sync"
)

var wg sync.WaitGroup

type node struct {
	word string
	level int
	parent string
}


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

func threading(c chan []string, word string, parentMap map[string]string) {
	defer wg.Done()
	var words []string
	for _, w := range get(buildURL(word)) {
		words = append(words, w)
		parentMap[w] = word
	}
	c <- words
}

func main() {
	word := "kiss"
	fmt.Println("START", word)
	maxDepth := 3

	//bfs
	var q map[string]int

	nq := map[string]int {
		word : 0,
	}

	vis := make(map[string]bool)

	store := make(map[node]bool)

	parentMap := make(map[string]string)

	for i := 1; i <= maxDepth; i++ {
		fmt.Println(i)
		q, nq = nq, make(map[string]int)
		queue := make(chan []string, 5000)
		for word := range q {
			if _, ok := vis[word]; !ok {
				wg.Add(1)
				vis[word] = true
				go threading(queue, word, parentMap)
			}
		}
		wg.Wait()
		close(queue)

		for v := range queue {
			for _, w := range v {
				nq[w] = i
				store[node { w, i, parentMap[w] }] = true
			}
		}
	}

	var ret []node

	ret = append(ret, node { word, 0, parentMap[word] })
	for k, _ := range store {
		ret = append(ret, k)
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].level < ret[j].level
	})

	for _, v := range ret {
		fmt.Println(v.word, " ", v.level, " ", v.parent)

	}

	fmt.Println("END")
}

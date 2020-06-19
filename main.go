package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	"net/http"
	"html/template"
	"sort"
	"strconv"
	"sync"

)

var wg sync.WaitGroup

type Node struct {
	Word string
	Level int
	Parent string
	Children string
}

type Page struct {
	Word string
	Nodes []Node
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

func threading(c chan []string, word string, parentMap, childrenMap map[string]string) {
	defer wg.Done()
	var words []string
	childString := ""
	for i, w := range get(buildURL(word)) {
		if i <= 1 {
			words = append(words, w)
			parentMap[w] = word
			childString += w + "/"
		}
		if i == 2 {
			words = append(words, w)
			parentMap[w] = word
			childString += w
			break
		}
  	}
	childrenMap[word] = childString
	c <- words
}

func wordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	word := vars["word"]
	fmt.Println("START", word)
	maxDepth, _ := strconv.Atoi(vars["depth"])

	//bfs
	var q map[string]int

	nq := map[string]int {
		word : 0,
	}

	vis := make(map[string]bool)

	store := make(map[Node]bool)

	parentMap := make(map[string]string)

	childrenMap := make(map[string]string)

	for i := 1; i <= maxDepth; i++ {
		fmt.Println(i)
		q, nq = nq, make(map[string]int)
		queue := make(chan []string, 5000)
		for word := range q {
			if _, ok := vis[word]; !ok {
				wg.Add(1)
				vis[word] = true
				go threading(queue, word, parentMap, childrenMap)
			}
		}
		wg.Wait()
		close(queue)
		for v := range queue {
			for _, w := range v {
				nq[w] = i
				store[Node { w, i, parentMap[w], "" }] = true
			}
		}
	}


	var List []Node

	List = append(List, Node { word, 0, parentMap[word], childrenMap[word] })
	for k, _ := range store {
		List = append(List, Node { k.Word, k.Level, k.Parent, childrenMap[k.Word]})
	}

	sort.Slice(List, func(i, j int) bool {
		return List[i].Level < List[j].Level
	})

	p := Page { Word: word, Nodes: List }

	t, _ := template.ParseFiles("Explore.html")
	t.Execute(w, p)

	fmt.Print(List)

	fmt.Println("END")
}


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/word/{word}/{depth}", wordHandler)
	http.ListenAndServe(":90", r)
}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

type Index struct {
	Path string `json:"path"`
}

func main() {
	client := &http.Client{}
	modules := make(map[string]int)
	versions := make(map[string]int)

	request, err := http.NewRequest("GET", "https://index.golang.org/index", nil)

	request.Header.Set("Disable-Module-Fetch", "true")

	if nil != err {
		log.Fatal(err)
	}

	response, err := client.Do(request)

	if nil != err {
		log.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode < 400 {
		s := bufio.NewScanner(response.Body)

		for s.Scan() {
			var index Index

			if err := json.Unmarshal(s.Bytes(), &index); nil != err {
				log.Fatal(err)
			}

			// Incrémente le nombre de modules et de versions pour la forge
			forge := strings.Split(index.Path, "/")[0]
			modules[forge]++
			versions[forge]++

		}

		// Crée un écrivain tabulaire
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 20, 4, 0, ' ', 0)
		defer w.Flush()

		// Écrit les colonnes du tableau
		fmt.Fprintf(w, "Forge\tModules\tVersions\n")

		// Parcours les forges triées par nombre de versions, décroissant
		for _, forge := range sortKeys(versions, func(a, b int) bool {
			return a > b

		}) {
			// Écrit les données pour chaque forge
			fmt.Fprintf(w, "%s\t%d\t%d\n", forge, modules[forge], versions[forge])
		}
	}
}

func sortKeys(m map[string]int, less func(a, b int) bool) []string {
	// Crée un tableau de clés
	keys := make([]string, len(m))

	// Parcours la map et ajoute les clés au tableau
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	// Trie le tableau par clés
	sort.Slice(keys, func(i, j int) bool {
		return less(m[keys[i]], m[keys[j]])
	})

	return keys

}

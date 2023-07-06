package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

func revip(s string) [][]string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://rapiddns.io/sameip/"+s+
		"?full=1#result", nil)
	res, err := client.Do(req)
	if err != nil {
	}
	defer res.Body.Close()

	re := regexp.MustCompile(`</th>\n<td>(.*?)</td>`)
	body, _ := ioutil.ReadAll(res.Body)
	r := re.FindAllStringSubmatch(string(body), -1)
	return r
}

func init() {
	flag.Usage = func() {
		h := []string{
			"",
			"revip (Reverse IP)",
			"",
			"By : github.com/vsec7",
			"",
			"Basic Usage :",
			" ▶ echo domain.com | revip",
			" ▶ cat listurls.txt | revip",
			"",
			"Options :",
			"  -o, --output        Output to file",
			"",
			"",
		}
		fmt.Fprintf(os.Stderr, strings.Join(h, "\n"))
	}
}

func main() {

	var outputFile string
	flag.StringVar(&outputFile, "output", "", "Output File")
	flag.StringVar(&outputFile, "o", "", "Output File")

	flag.Parse()

	jobs := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i <= 10; i++ {
		wg.Add(1)

		go func() {
			for u := range jobs {

				ri := revip(u)
				for _, d := range ri {
					fmt.Println(d[1])

					if len(outputFile) != 0 {
						file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							fmt.Printf("[!] Failed Creating File: %s", err)
						}
						buf := bufio.NewWriter(file)
						buf.WriteString(d[1] + "\n")
						buf.Flush()
						file.Close()
					}
				}
			}
			wg.Done()
		}()
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		jobs <- sc.Text()
	}
	close(jobs)
	wg.Wait()
}

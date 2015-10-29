package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"bufio"
	"flag"
	"log"
	"fmt"
	"os"
)

func readDictFile(file string) (map[string]string, error) {
	res := map[string]string{}
	f, err := os.Open(file)
	if err != nil {
		return res, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\n")
		elems := strings.SplitN(line, "\t", 2)
		res[elems[0]] = elems[1]
	}
	if err := scanner.Err(); err != nil {
		return res, err
	}
	return res, nil
}

func writeDictFile(file string, dict map[string]string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	for word, definition := range dict {
		_, err = f.WriteString(fmt.Sprintf("%s\t%s\n", word, definition))
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	baseUrl := flag.String("url", "http://www.jahanshiri.ir/pvc/conjugate?script=p&level=d", "URL of the service to query")
	lang := flag.String("lang", "fa", "Language to use (en, fa, fr, ...)")
	wordlist := flag.String("wordlist", "words.txt", "File containing the list of word (one per line) to query")
	dictfile := flag.String("dict", "persian-verbs.txt", "Dictionary file (format is \"word\\tdefinition\\n\"...)")
	replace := flag.Bool("replace", false, "Force to query again existing words in dictionary")
	flag.Parse()

	dict, err := readDictFile(*dictfile)
	if err != nil {
		if os.IsNotExist(err) {
			dict = map[string]string{}
		} else {
			log.Fatal(err)
		}
	}

	u, err := url.Parse(*baseUrl)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(*wordlist)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		word := strings.TrimRight(scanner.Text(), "\n")
		if len(word) == 0 {
			continue
		}
		if _, has_word := dict[word]; !*replace && has_word {
			fmt.Printf("Skip word: %s\n", word)
			continue
		}
		params := u.Query()
		params.Set("verb", word)
		params.Set("lang", *lang)
		params.Set("passive", "no")
		u.RawQuery = params.Encode()
		res, err := http.Get(u.String())
		if err != nil {
			log.Println(err)
			continue
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			continue
		}
		params.Set("passive", "no")
		u.RawQuery = params.Encode()
		res, err = http.Get(u.String())
		if err != nil {
			log.Println(err)
			continue
		}
		bodyPassive, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			continue
		}
		definition := fmt.Sprintf(html, css, string(body), string(bodyPassive))
		definition = strings.Replace(definition, `\`, `\\`, -1)
		definition = strings.Replace(definition, "\n", `\n`, -1)
		dict[word] = definition
		err = writeDictFile(*dictfile, dict)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf("Saved word: %s\n", word)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

const html string = `<!DOCTYPE html>
<html><head><style type="text/css">
%s
</style></head><body>
%s
<hr/>
%s
</body></html>`

const css string = `
/* http://www.jahanshiri.ir/pvc/pvc.css */

#btnConj { background-position:-1px -1px; }
.rtl #btnConj { background-position:-32px -1px; }

.tblConj table { display:inline-block; margin:1em .5em; }
.tblConj-rtl table, .tblConj-rtl th  { text-align: right; direction: rtl; }
.tblConj-ltr table, .tblConj-ltr th  { text-align: left; direction: left; }

#top-verbs { display:none; }

/* http://www.jahanshiri.ir/inc/common.css */

table {
	border:0;
	border-collapse:collapse;
	background-color:transparent;
	font-size:100%;
	margin:1em 0
}
caption {
	text-align:center;
	padding:2px;
	white-space:nowrap
}
tr {
	vertical-align:middle
}
.va-top tr,.va-top {
	vertical-align:top
}
th,td {
	padding:4px 8px;
	border:1px solid #c4b7a9
}
th[colspan],td[colspan] {
	text-align:center
}
th {
	background-color:#f1e5d8;
	font-weight:normal;
	text-align:left
}
.thh,.thh th {
	background-color:#e8d4be
}
.td-ltr caption {
	direction:ltr
}
.rtl th {
	text-align:right
}
`


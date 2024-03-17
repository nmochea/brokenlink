package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	thread *int
	silent *bool
	ua     *string
)

func req(url string) {
	if !strings.Contains(url, "http") {
		fmt.Println("\033[31m[-]\033[37m Send URLs via stdin (ex: cat js.txt | Brokenlink)")
		os.Exit(0)
	}
	var secretslist = []string{"https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)","http?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)"}

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	transp := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpclient := &http.Client{Transport: transp}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", *ua)
	r, err := httpclient.Do(req)
	if err != nil {
		fmt.Println("\033[31m[-]\033[37m", "\033[37m"+"Unable to make a request for " + url + "\033[37m")
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("\033[31m[-]\033[37m", "\033[37m"+"Unable to read the body of " + url + "\033[37m")
	}
	strbody := string(body)

	for _, secret := range secretslist {
		match, err := regexp.MatchString(secret, strbody)
		if err != nil {
			fmt.Println("\033[31m[-]\033[37m", "\033[37m"+"Regex Error " + url + "\033[37m")
		}
		if match {
			pattern := regexp.MustCompile(secret)
			matched := pattern.FindString(strbody)
			fmt.Println("\033[32m[+]\033[37m", url, "\033[32m[\033[37m" + matched + "\033[32m]\033[37m")
		}
	}
}

func init() {
	silent = flag.Bool("s", false, "silent")
	thread = flag.Int("t", 50, "thread number")
	ua = flag.String("ua", "Brokenlink", "User-Agent")
}

func banner() {
	fmt.Println("\033[31m" + `
	 ___         _            _ _      _   
 	| _ )_ _ ___| |_____ _ _ | (_)_ _ | |__
 	| _ \ '_/ _ \ / / -_) ' \| | | ' \| / /
 	|___/_| \___/_\_\___|_||_|_|_|_||_|_\_\
	` + "\033[31m[\033[37mNmochea v0.1\033[31m]\n")
}

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	urls := make(chan string)
	var wg sync.WaitGroup
	flag.Parse()

	if !*silent {
		banner()
		for i := 0; i < *thread; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for url := range urls {
					req(url)
				}
			}()
		}
		for stdin.Scan() {
			url := stdin.Text()
			urls <- url
		}
		close(urls)
		wg.Wait()
		//os.Exit(1)
	} else {
		for i := 0; i < *thread; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for url := range urls {
					req(url)
				}
			}()
		}
		for stdin.Scan() {
			url := stdin.Text()
			urls <- url
		}
		close(urls)
		wg.Wait()
		//os.Exit(1)
	}
}

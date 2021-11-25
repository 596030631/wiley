package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"math/rand"
	"net"
	"net/http"
)

const seleniumPath = "./chromedriver.exe"
const Core = 5

var param = make([]chan string, Core)
var result = make([]chan string, Core)

func indexHttp(w http.ResponseWriter, r *http.Request) {
	k := r.FormValue("keyword")
	p := r.FormValue("page")
	if k != "" && p != "" {
		id := rand.Intn(Core)
		fmt.Printf("keyword=%s page=%s\n", k, p)
		go func() {
			param[id] <- fmt.Sprintf("https://onlinelibrary.wiley.com/action/doSearch?Ppub=&field1=Keyword&text1=%s&pageSize=20&startPage=%s", k, p)
		}()
		_, _ = w.Write([]byte(<-result[id]))
	} else {
		_, _ = w.Write([]byte("param null"))
	}
}

func startSelenium(id int) {

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return
	}

	selenium.SetDebug(false)
	_, _ = selenium.NewChromeDriverService(seleniumPath, port)

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	imgCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imgCaps,
		Path:  "",
		Args: []string{
			"--headless",
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36",
		},
	}
	caps.AddChrome(chromeCaps)

	browner, e := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))

	if e != nil {
		fmt.Println(e)
	}

	for {
		url := <-param[id]
		_ = browner.Get(url)
		html, _ := browner.PageSource()
		result[id] <- html
	}
}

func main() {
	for i := 0; i < Core; i++ {
		param[i] = make(chan string, Core)
		result[i] = make(chan string, Core)
		go startSelenium(i)
	}
	http.HandleFunc("/wiley", indexHttp)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

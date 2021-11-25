package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"net"
	"net/http"
)

const seleniumPath = "./chromedriver.exe"

var param = make(chan string)
var result = make(chan string)

func indexHttp(w http.ResponseWriter, r1 *http.Request) {
	k := r1.FormValue("keyworld")
	p := r1.FormValue("page")
	if k != "" && p != "" {
		param <- fmt.Sprintf("https://onlinelibrary.wiley.com/action/doSearch?Ppub=&field1=Keyword&text1=%s&pageSize=20&startPage=%s")
		_, _ = w.Write([]byte(<-result))
		fmt.Println("返回网页结果")
	} else {
		_, _ = w.Write([]byte("参数不全"))
	}
}

func startSelenium() {

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
	service, _ := selenium.NewChromeDriverService(seleniumPath, port)
	defer func(service *selenium.Service) {
		err := service.Stop()
		if err != nil {
			return
		}
	}(service)
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
			//"--headless",
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36",
		},
	}
	caps.AddChrome(chromeCaps)

	browner, e := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port+i))
	if err != nil {
		return
	}
	if e != nil {
		fmt.Println(e)
	}

	var window = [3]string{}

	for i := 0; i < 3; i++ {
		script, _ := browner.ExecuteScript(fmt.Sprintf(`window.open("%s", "_blank");`, url), nil)
	}

	for {
		url := <-param
		script, _ := browner.ExecuteScript(fmt.Sprintf(`window.open("%s", "_blank");`, url), nil)
		_ = browner.Get(url)

		//print(script)
		html, _ := browner.PageSource()
		browner.SwitchSession(session)
		browner.Get()
		result <- html
	}
}

func main() {

	go startSelenium()

	http.HandleFunc("/wiley", indexHttp)

	log.Fatal(http.ListenAndServe(":8081", nil))

}

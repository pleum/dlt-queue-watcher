package main

import (
	"encoding/json"
	"log"
	"net/http/cookiejar"
	"strings"

	"github.com/go-resty/resty/v2"
)

var id = ""
var pass = ""

func main() {
	client := resty.New()

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(15))
	client.SetCookieJar(jar)

	resp, err := client.R().Get("https://gecc.dlt.go.th/")
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	a := resp.RawResponse.Request.URL.String()

	re, err := client.R().SetFormData(map[string]string{
		"per_id":       id,
		"per_password": pass,
	}).Post(a + "checklogin.php")
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	isPasswordInvalid := strings.Contains(re.String(), "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง")
	if isPasswordInvalid {
		return
	}

	resa, err := client.R().Get(a + "home.php?selectProvice=3&selectOffice=12&selectGroup=null&selectWork=105921")
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	bo := resa.String()

	t := "var my_events ="
	i := strings.Index(bo, t)

	c := bo[i+len(t)+1:]

	n := strings.Index(c, ";")

	g := c[:n]

	var q []Queue

	if err := json.Unmarshal([]byte(g), &q); err != nil {
		panic(err)
	}
	for _, val := range q {
		log.Println(val.Title, val.Start)
	}
}

type Queue struct {
	Title string `json:"title"`
	Start string `json:"start"`
}

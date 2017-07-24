package main

import (
	"bufio"
	"flag"
	"fmt"
	"path"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/zc310/utils"
	"github.com/zc310/utils/httputil"

	"os"
	"regexp"
	"strconv"

	"io/ioutil"
)

func main() {
	os.MkdirAll("img", 0777)
	var Bid int
	var err2 error
	tid := flag.Int("id", 0, "商品编号")
	flag.Parse()
	if *tid == 0 {
		goto NID
	} else {
		Bid = *tid
	}
NID:
	fmt.Printf("\n请输入商品编号）: ")
	reader := bufio.NewReader(os.Stdin)
	line, _, _ := reader.ReadLine()

	re := regexp.MustCompile(`\d+`)
	for _, t := range re.FindAllString(string(line), -1) {
		Bid, err2 = strconv.Atoi(t)
		if err2 != nil {
			continue
		}
		if err := down(Bid); err != nil {
			fmt.Println(err.Error())
		}
	}

	goto NID
}
func saveimg(savepath string, url string) {
	fmt.Println(url)
	b, err := httputil.GetByte(url)
	if err != nil {
		return
	}

	ioutil.WriteFile(fmt.Sprintf("%s%s%s", savepath, utils.MD5([]byte(url)), path.Ext(url)), b, 0777)
}
func down(id int) error {
	res, err := httputil.GetRequest(fmt.Sprintf("http://item.jd.com/%d.html", id))
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`.com/n(\d)/`)

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("img/%s/%d/", time.Now().Format("2006-01-02"), id)
	os.MkdirAll(path, 0777)

	doc.Find(".spec-items img").Each(func(i int, s *goquery.Selection) {
		if url, ok := s.Attr("src"); ok {
			url := re.ReplaceAllString(url, ".com/n1/")
			if len(url) > 0 {
				saveimg(path, "http:"+url)
			}

		}

	})


	s, err := httputil.GetString(fmt.Sprintf("http://d.3.cn/desc/%d?cdn=2&callback=showdesc", id))

	if err != nil {
		return err
	}

	re = regexp.MustCompile(`//[\w\-\.\/]+.jpg`)

	for _, url := range re.FindAllString(s, -1) {
		saveimg(path, "http:"+url)
	}

	return nil

}

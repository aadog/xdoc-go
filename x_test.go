package xdoc

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"testing"
)

func TestDoc(t *testing.T) {
	d, err := LoadDoc("x.html")
	if err != nil {
		panic(err)
	}

	e := d.FindOne(`//div[@class="formArea1 formAreaCom"]`).CountChild()
	fmt.Println(e)
}

func TestGoQuery(t *testing.T) {
	f, err := os.Open("x.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(doc.Find("xxzzzzasdf"))
}

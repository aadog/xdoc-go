package xdoc

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"os"
)

type Node struct {
	node *html.Node
}

func (n *Node) InnerText() string {
	return htmlquery.InnerText(n.node)
}
func (n *Node) OutputHTML(self bool) string {
	return htmlquery.OutputHTML(n.node, self)
}
func (n *Node) ExistsAttr(name string) bool {
	return htmlquery.ExistsAttr(n.node, name)
}
func (n *Node) Attr(name string) string {
	return htmlquery.SelectAttr(n.node, name)
}
func (n *Node) AttrOr(name string, val string) string {
	if n.ExistsAttr(name) && n.Attr(name) != "" {
		return n.Attr(name)
	}
	return val
}
func (n *Node) IsNotFind() bool {
	return n.node.Type == html.ErrorNode
}
func (n *Node) QueryOne(expr string) (*Node, error) {
	newNode, err := htmlquery.Query(n.node, expr)
	if err != nil {
		return nil, err
	}
	newN := NodeFromHtmlNode(newNode)
	if newN.IsNotFind() == true {
		return nil, errors.New(fmt.Sprintf("not found xpath:%s", expr))
	}
	return newN, nil
}
func (n *Node) FindOne(expr string) *Node {
	newN := NodeFromHtmlNode(htmlquery.FindOne(n.node, expr))
	return newN
}
func (n *Node) FindAll(expr string) []*Node {
	newNs := make([]*Node, 0)
	ds := htmlquery.Find(n.node, expr)
	for _, d := range ds {
		newN := NodeFromHtmlNode(d)
		newNs = append(newNs, newN)
	}
	return newNs
}

func (n *Node) Tag() string {
	return n.node.Data
}
func (n *Node) Children() []*Node {
	return n.FindAll("./*")
}
func (n *Node) FirstChild() *Node {
	return n.FindOne(`./*[1]`)
}
func (n *Node) LastChild() *Node {
	return n.FindOne(`./*[last()]`)
}

//	func (n *Node) FindAll(expr string) []*Node {
//		newNs := make([]*Node, 0)
//
//		for _, d := range ds {
//			newN := &Node{}
//			newN.node = d
//			newNs = append(newNs, newN)
//		}
//		return newNs
//	}
func NodeFromHtmlNode(n *html.Node) *Node {
	newN := &Node{}
	newN.node = n
	if newN.node == nil {
		newN.node = &html.Node{}
	}
	return newN
}
func NewDocumentFromReader(r io.Reader) (*Node, error) {
	h, err := htmlquery.Parse(r)
	if err != nil {
		return nil, err
	}
	d := &Node{node: h}
	return d, nil
}
func LoadDoc(path string) (*Node, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewDocumentFromReader(bufio.NewReader(f))
}
func LoadURL(url string) (*Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	return NewDocumentFromReader(bufio.NewReader(r))
}

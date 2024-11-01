package link

import (
	"io"
	"strings"

	stack "github.com/ntsiris/sitemap-builder/pkg/collections"
	"golang.org/x/net/html"
)

type HTMLLink struct {
	Href string
	Text string
}

type nodeProcessor interface {
	qualifies(node *html.Node) bool
	processNode(node *html.Node)
}

type textNodeProcessor struct {
	textBuilder *strings.Builder
}

func (tp *textNodeProcessor) qualifies(node *html.Node) bool {
	return node.Type == html.TextNode
}

func (tp *textNodeProcessor) processNode(node *html.Node) {
	tp.textBuilder.WriteString(strings.TrimSpace(node.Data))
}

type linkNodeProcessor struct{}

func (lp *linkNodeProcessor) qualifies(node *html.Node) bool {
	return node.Type == html.ElementNode && node.Data == "a"
}

func (lp *linkNodeProcessor) processNode(node *html.Node) {}

func Parse(r io.Reader) ([]HTMLLink, error) {
	htmlDocument, err := html.Parse(r)

	if err != nil {
		return nil, err
	}

	nodes := dfsHTMLNodes(htmlDocument, &linkNodeProcessor{})

	var links []HTMLLink

	for _, node := range nodes {
		links = append(links, buildHTMLLink(node))
	}

	return links, nil
}

func dfsHTMLNodes(seed *html.Node, pr nodeProcessor) []*html.Node {
	var ret []*html.Node
	visited := make(map[*html.Node]bool)
	nodeStack := stack.Stack{}

	nodeStack.Push(seed)

	for nodeStack.Len() > 0 {
		node := nodeStack.Pop().(*html.Node)

		if !visited[node] {
			visited[node] = true

			for child := node.FirstChild; child != nil; child = child.NextSibling {
				nodeStack.Push(child)
			}

			if pr.qualifies(node) {
				pr.processNode(node)
				ret = append(ret, node)
			}
		}
	}

	return ret
}

func buildHTMLLink(node *html.Node) HTMLLink {
	var ret HTMLLink

	for _, attr := range node.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}

	ret.Text = getLinkText(node)

	return ret
}

func getLinkText(node *html.Node) string {

	tp := textNodeProcessor{}
	tp.textBuilder = &strings.Builder{}

	_ = dfsHTMLNodes(node, &tp)

	return tp.textBuilder.String()
}

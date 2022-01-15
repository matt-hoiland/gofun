package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const MarkdownFile = "/Users/matthew.hoiland/journal/albany/2022-01-14_Fr_daily.md"

func main() {

	fp, err := os.Open(MarkdownFile)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	data, err := io.ReadAll(fp)
	if err != nil {
		panic(err)
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)
	var buf bytes.Buffer
	if err := md.Convert(data, &buf); err != nil {
		panic(err)
	}
	data, err = io.ReadAll(&buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	doc := bytes.NewBuffer(data)
	root, err := html.Parse(doc)
	if err != nil {
		panic(err)
	}
	TraverseDOM(root, 0)

	var tasks []Task
	ExtractTasks(&tasks, root)
	for _, task := range tasks {
		fmt.Println(task)
	}
}

func TraverseDOM(node *html.Node, depth int) {
	if node == nil {
		return
	}
	indent := strings.Repeat("  ", depth)
	switch node.Type {
	case html.ElementNode:
		fmt.Printf("%s<%s(%d)", indent, node.DataAtom.String(), node.DataAtom)
		for _, attr := range node.Attr {
			fmt.Printf(" %s:%s=\"%s\"", attr.Namespace, attr.Key, attr.Val)
		}
		fmt.Println(">")
	case html.TextNode:
		text := strings.TrimSpace(node.Data)
		if text != "" {
			fmt.Printf("%s\"%s\"\n", indent, strings.TrimSpace(node.Data))
		}
	case html.CommentNode:
		fmt.Printf("%sCOMMENT: \"%s\"\n", indent, node.Data)
	case html.DoctypeNode:
		fmt.Printf("%sDOCTYPE: \"%s\"\n", indent, node.Data)
	case html.DocumentNode:
		fmt.Printf("%sDOCUMENT: \"%s\"\n", indent, node.Data)
	case html.ErrorNode:
		fmt.Printf("%sERROR: \"%s\"\n", indent, node.Data)
	case html.RawNode:
		fmt.Printf("%sRAW: \"%s\"\n", indent, node.Data)
	}

	// Depth first, traverse children
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		TraverseDOM(child, depth+1)
	}
}

type Task struct {
	Completed bool
	Name      string
}

func (t Task) String() string {
	tern := func(b bool, t string, f string) string {
		if b {
			return t
		}
		return f
	}
	return fmt.Sprintf("[%s] %s", tern(t.Completed, "x", " "), t.Name)
}

func ExtractTasks(tasks *[]Task, node *html.Node) {
	if node == nil {
		return
	}
	if node.Type == html.ElementNode && node.DataAtom == atom.Li {
		var input, text *html.Node
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode && child.DataAtom == atom.Input {
				input = child
			}
			if child.Type == html.TextNode {
				text = child
			}
		}
		if input == nil || text == nil {
			return
		}

		isCheckBox := false
		completed := false
		for _, attr := range input.Attr {
			isCheckBox = isCheckBox || attr.Key == "type" && attr.Val == "checkbox"
			completed = completed || attr.Key == "checked"
		}
		if isCheckBox {
			*tasks = append(*tasks, Task{completed, strings.TrimSpace(text.Data)})
		}
		return
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		ExtractTasks(tasks, c)
	}
}

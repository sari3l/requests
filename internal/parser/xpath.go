package parser

import (
	"bytes"
	"fmt"
	AXpath "github.com/antchfx/xpath"
	"github.com/golang/groupcache/lru"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"log"
	"strings"
	"sync"
)

type XpathNode html.Node

// 基于https://github.com/antchfx/htmlquery，但切换更方便的处理方式
var _ AXpath.NodeNavigator = &NodeNavigator{}

func CreateXPathNavigator(top *html.Node) *NodeNavigator {
	return &NodeNavigator{curr: top, root: top, attr: -1}
}

func XpathParser(body *string) *XpathNode {
	reader := strings.NewReader(*body)
	node, err := html.Parse(reader)
	if err != nil {
		log.Printf("%+v\n", errors.WithStack(err))
		return nil
	}
	return (*XpathNode)(node)
}

func (x *XpathNode) Find(expr string) []*XpathNode {
	nodes, err := x.QueryAll(expr)
	if err != nil {
		panic(err)
	}
	return nodes
}

func (x *XpathNode) FindOne(expr string) *XpathNode {
	node, err := x.Query(expr)
	if err != nil {
		panic(err)
	}
	return node
}

func (x *XpathNode) QueryAll(expr string) ([]*XpathNode, error) {
	exp, err := getQuery(expr)
	if err != nil {
		return nil, err
	}
	nodes := x.QuerySelectorAll(exp)
	return nodes, nil
}

func (x *XpathNode) Query(expr string) (*XpathNode, error) {
	exp, err := getQuery(expr)
	if err != nil {
		return nil, err
	}
	return x.QuerySelector(exp), nil
}

// QuerySelector returns the first matched html.Node by the specified XPath selector.
func (x *XpathNode) QuerySelector(selector *AXpath.Expr) *XpathNode {
	t := selector.Select(CreateXPathNavigator((*html.Node)(x)))
	if t.MoveNext() {
		return getCurrentNode(t.Current().(*NodeNavigator))
	}
	return nil
}

func (x *XpathNode) QuerySelectorAll(selector *AXpath.Expr) []*XpathNode {
	var elems []*XpathNode
	t := selector.Select(CreateXPathNavigator((*html.Node)(x)))
	for t.MoveNext() {
		nav := t.Current().(*NodeNavigator)
		n := getCurrentNode(nav)
		elems = append(elems, n)
	}
	return elems
}

func getCurrentNode(n *NodeNavigator) *XpathNode {
	if n.NodeType() == AXpath.AttributeNode {
		childNode := &html.Node{
			Type: html.TextNode,
			Data: n.Value(),
		}
		return &XpathNode{
			Type:       html.ElementNode,
			Data:       n.LocalName(),
			FirstChild: childNode,
			LastChild:  childNode,
		}

	}
	return (*XpathNode)(n.curr)
}

// functions

func (x *XpathNode) Text() string {
	var output func(*bytes.Buffer, *html.Node)
	output = func(buf *bytes.Buffer, n *html.Node) {
		switch n.Type {
		case html.TextNode:
			buf.WriteString(n.Data)
			return
		case html.CommentNode:
			return
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			output(buf, child)
		}
	}

	var buf bytes.Buffer
	output(&buf, (*html.Node)(x))
	return buf.String()
}
func (x *XpathNode) SelectAttr(name string) (val string) {
	if x == nil {
		return
	}
	if x.Type == html.ElementNode && x.Parent == nil && name == x.Data {
		return x.Text()
	}
	for _, attr := range x.Attr {
		if attr.Key == name {
			val = attr.Val
			break
		}
	}
	return
}

// ExistsAttr returns whether attribute with specified name exists.
func (x *XpathNode) ExistsAttr(name string) bool {
	if x == nil {
		return false
	}
	for _, attr := range x.Attr {
		if attr.Key == name {
			return true
		}
	}
	return false
}

// HTML returns the text including tags name.
func (x *XpathNode) HTML(self bool) string {
	var buf bytes.Buffer
	if self {
		html.Render(&buf, (*html.Node)(x))
	} else {
		for x := x.FirstChild; x != nil; x = x.NextSibling {
			html.Render(&buf, x)
		}
	}
	return buf.String()
}

//

type NodeNavigator struct {
	root, curr *html.Node
	attr       int
}

func (h *NodeNavigator) Current() *html.Node {
	return h.curr
}

func (h *NodeNavigator) NodeType() AXpath.NodeType {
	switch h.curr.Type {
	case html.CommentNode:
		return AXpath.CommentNode
	case html.TextNode:
		return AXpath.TextNode
	case html.DocumentNode:
		return AXpath.RootNode
	case html.ElementNode:
		if h.attr != -1 {
			return AXpath.AttributeNode
		}
		return AXpath.ElementNode
	case html.DoctypeNode:
		// ignored <!DOCTYPE Html> declare and as Root-XpathNode type.
		return AXpath.RootNode
	}
	panic(fmt.Sprintf("unknown Html node type: %v", h.curr.Type))
}

func (h *NodeNavigator) LocalName() string {
	if h.attr != -1 {
		return h.curr.Attr[h.attr].Key
	}
	return h.curr.Data
}

func (*NodeNavigator) Prefix() string {
	return ""
}

func (h *NodeNavigator) Value() string {
	switch h.curr.Type {
	case html.CommentNode:
		return h.curr.Data
	case html.ElementNode:
		if h.attr != -1 {
			return h.curr.Attr[h.attr].Val
		}
		return (*XpathNode)(h.curr).Text()
	case html.TextNode:
		return h.curr.Data
	}
	return ""
}

func (h *NodeNavigator) Copy() AXpath.NodeNavigator {
	n := *h
	return &n
}

func (h *NodeNavigator) MoveToRoot() {
	h.curr = h.root
}

func (h *NodeNavigator) MoveToParent() bool {
	if h.attr != -1 {
		h.attr = -1
		return true
	} else if node := h.curr.Parent; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *NodeNavigator) MoveToNextAttribute() bool {
	if h.attr >= len(h.curr.Attr)-1 {
		return false
	}
	h.attr++
	return true
}

func (h *NodeNavigator) MoveToChild() bool {
	if h.attr != -1 {
		return false
	}
	if node := h.curr.FirstChild; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *NodeNavigator) MoveToFirst() bool {
	if h.attr != -1 || h.curr.PrevSibling == nil {
		return false
	}
	for {
		node := h.curr.PrevSibling
		if node == nil {
			break
		}
		h.curr = node
	}
	return true
}

func (h *NodeNavigator) String() string {
	return h.Value()
}

func (h *NodeNavigator) MoveToNext() bool {
	if h.attr != -1 {
		return false
	}
	if node := h.curr.NextSibling; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *NodeNavigator) MoveToPrevious() bool {
	if h.attr != -1 {
		return false
	}
	if node := h.curr.PrevSibling; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *NodeNavigator) MoveTo(other AXpath.NodeNavigator) bool {
	node, ok := other.(*NodeNavigator)
	if !ok || node.root != h.root {
		return false
	}

	h.curr = node.curr
	h.attr = node.attr
	return true
}

// cache

var DisableSelectorCache = false
var SelectorCacheMaxEntries = 50

var (
	cacheOnce  sync.Once
	cache      *lru.Cache
	cacheMutex sync.Mutex
)

func getQuery(expr string) (*AXpath.Expr, error) {
	if DisableSelectorCache || SelectorCacheMaxEntries <= 0 {
		return AXpath.Compile(expr)
	}
	cacheOnce.Do(func() {
		cache = lru.New(SelectorCacheMaxEntries)
	})
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if v, ok := cache.Get(expr); ok {
		return v.(*AXpath.Expr), nil
	}
	v, err := AXpath.Compile(expr)
	if err != nil {
		return nil, err
	}
	cache.Add(expr, v)
	return v, nil
}

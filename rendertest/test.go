package main

import (
	"log"
	"os"
	"strings"

	"github.com/thomgray/codebook/htmlrender"
	"github.com/thomgray/egg"
	"golang.org/x/net/html"
)

func main() {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	defer file.Close()

	htmls := `<h1>hello</h1>
		<p>
		one two three four <code>five</code> six seven <code>eight</code> nine ten <code>eleven</code>
		twelve <a href="https://www.isthisworking.com">thirteen</a>
		</p>
		`

	node, _ := html.Parse(strings.NewReader(htmls))

	app := egg.InitOrPanic()
	defer app.Start()

	app.OnDraw(func(c egg.Canvas) {
		htmlrender.RenderHtml(node, c)
	})
}

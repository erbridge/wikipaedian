package wiki

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"cgt.name/pkg/go-mwclient"
	"github.com/erbridge/gotwit"
)

type (
	Client struct {
		bot    *gotwit.Bot
		wiki   *mwclient.Client
		ticker *time.Ticker
	}
)

func NewClient(bot *gotwit.Bot) (client Client, err error) {
	// TODO: Add contact details to the user agent.
	wiki, err := mwclient.New("https://en.wikipedia.org/w/api.php", "Wikipaedian")

	if err != nil {
		return
	}

	client = Client{
		bot:  bot,
		wiki: wiki,
	}

	return
}

func (c *Client) Start(d time.Duration) {
	c.ticker = time.NewTicker(d)
	defer c.ticker.Stop()

	c.Post()

	for range c.ticker.C {
		c.Post()
	}
}

func (c *Client) Post() {
	last := c.LastPost()

	if last == "" {
		c.bot.Post("[[\"Hello, World!\" program|Hello World!]]", false)

		return
	}

	page, _ := c.Page(last)

	content := c.CreatePost(page)

	fmt.Println(content)

	// c.bot.Post(content, false)
}

func (c *Client) LastPost() string {
	if text, err := c.bot.LastTweetText(); err != nil {
		panic(err)
	} else {
		return text
	}

	return ""
}

func (c *Client) Page(last string) (content, timestamp string) {
	re := regexp.MustCompile("\\[{2}([^:]+?)\\]{2}")

	matches := re.FindAllStringSubmatch(last, -1)

	// FIXME: Handle the 0 match state.
	title := matches[rand.Intn(len(matches))][1]

	re = regexp.MustCompile("(.+?)\\|")

	match := re.FindStringSubmatch(title)

	if len(match) > 0 {
		title = match[1]
	}

	content, timestamp, err := c.wiki.GetPageByName(title)

	if err != nil {
		panic(err)
	}

	return
}

func (c *Client) CreatePost(page string) (content string) {
	re := regexp.MustCompile("[.!?][[:space:]]+" +
		"(" +
		"(" +
		"([^{^}^\\[^\\]^<^>]+?[[:space:]]+?)+?" +
		")" +
		"(\\[{2}[^:]+?\\]{2})" +
		"(" +
		"([[:space:]]+?[^{^}^\\[^\\]^<^>]+?)+?" +
		")" +
		"([.!?])" +
		")" +
		"[[:space:]]+",
	)

	matches := re.FindAllStringSubmatch(page, -1)

	// FIXME: Handle the 0 match state.
	match := matches[rand.Intn(len(matches))]

	content = match[1]

	if len(content) <= 140 {
		return
	}

	content = match[4]

	if len(content) > 138 {
		re = regexp.MustCompile("\\[{2}[^:]+?\\]{2}")

		matches := re.FindAllString(content, -1)

		// FIXME: Handle the 0 match state.
		content = matches[rand.Intn(len(matches))]
	}

	before, after := strings.Split(match[2], " "), strings.Split(match[5], " ")

	newContent := content

	for len(newContent) < 138 {
		content = newContent

		word := ""

		if len(before) > 0 && rand.Float32() > 0.5 {
			word, before = before[len(before)-1], before[:len(before)-1]

			if word == "" {
				continue
			}

			newContent = word + " " + content
		} else if len(after) > 0 {
			word, after = after[0], after[1:]

			if word == "" {
				continue
			}

			newContent = content + " " + word
		} else {
			break
		}
	}

	if len(before) > 0 {
		content = "…" + content
	}

	if len(after) > 0 {
		content += "…"
	} else {
		content += match[7]
	}

	return
}

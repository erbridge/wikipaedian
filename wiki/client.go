package wiki

import (
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

func NewClient(b gotwit.Bot) (client Client, err error) {
	// TODO: Add contact details to the user agent.
	wiki, err := mwclient.New("https://en.wikipedia.org/w/api.php", "Wikipaedian")

	if err != nil {
		return
	}

	client = Client{
		wiki: wiki,
	}

	return
}

func (c *Client) Start(d time.Duration) {
	c.ticker = time.NewTicker(d)
	defer c.ticker.Stop()

	println("Tick")

	for range c.ticker.C {
		println("Tick")
	}
}

package fluent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Client to send data to fluentd.
type Client struct {
	buff chan map[string]interface{}
	addr string
	envs map[string]string
}

// NewClient creates client and runs worker.
func NewClient(host string, port int, tag string) *Client {
	c := &Client{
		buff: make(chan map[string]interface{}, 1024),
		addr: fmt.Sprintf("http://%v:%v/%v", host, port, tag),
		envs: map[string]string{},
	}

	if len(os.Args) >= 1 {
		c.envs["appname"] = os.Args[0]
	}

	c.envs["podname"], _ = os.Hostname()

	go c.work()

	return c
}

// Push pushes data to buffer.
func (c *Client) Push(data map[string]interface{}) {
	select {
	case c.buff <- data:
	default:
	}
}

func (c *Client) work() {
	ticker := time.NewTicker(time.Second)
	var arr []map[string]interface{}

	for {
		select {
		case data := <-c.buff:
			for k, v := range c.envs {
				data[k] = v
			}

			arr = append(arr, data)
			if len(arr) <= 87 { // buff size
				continue
			}

		case <-ticker.C:
		}

		// do flush.
		if len(arr) == 0 {
			continue
		}

		b, err := json.Marshal(arr)
		arr = nil

		if err != nil {
			log.Println(err)
			continue
		}

		_, err = http.Post(c.addr, "application/json", bytes.NewBuffer(b))
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

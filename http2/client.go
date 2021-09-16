package http2

import (
	"bufio"
	"bytes"
	"fmt"
	// "crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	// "golang.org/x/net/http2"
)

type Client struct {
	client *http.Client
}

func (c *Client) Dial() {
	// Adds TLS cert-key pair
	// certs, err := tls.LoadX509KeyPair("./http2/certs/key.crt", "./http2/certs/key.key")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// t := &http2.Transport{
	// 	TLSClientConfig: &tls.Config{
	// 		Certificates:       []tls.Certificate{certs},
	// 		InsecureSkipVerify: true,
	// 	},
	// }
	// c.client = &http.Client{Transport: t}
	c.client = &http.Client{}
}

func (c *Client) Post(data []byte) {
	host := os.Getenv("POST_HOST")
	path := os.Getenv("POST_PATH")
	auth := fmt.Sprintf("Bearer %s", os.Getenv("POST_KEY"))
	
	log.Println("POST", host, path)
	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   host,
			Path:   path,
		},
		Header: http.Header{
			"Content-Type": []string{"application/json"},
			"Authorization": []string{auth},
		},
		Body:   ioutil.NopCloser(bytes.NewReader(data)),
	}
	// Sends the request
	start := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode == 500 {
		return
	}

	defer resp.Body.Close()

	bufferedReader := bufio.NewReader(resp.Body)

	buffer := make([]byte, 4*1024)

	var totalBytesReceived int

	// Reads the response
	for {
		len, err := bufferedReader.Read(buffer)
		if len > 0 {
			duration := time.Since(start)
			totalBytesReceived += len
			log.Println(duration, "elapsed", len, "bytes received")
			// Prints received data
			log.Println(string(buffer[:len]))
		}

		if err != nil {
			if err == io.EOF {
				// Last chunk received
				// log.Println(err)
			}
			break
		}
	}
	duration := time.Since(start)
	log.Println("Total Time Elapsed:", duration)
	log.Println("Total Bytes Sent:", len(data))
	log.Println("Total Bytes Received:", totalBytesReceived)
}

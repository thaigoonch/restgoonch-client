package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/golang/protobuf/proto"
	restgoonch "github.com/thaigoonch/restgoonch/service"
)

var (
	host = "restgoonch-service"
	port = 8080
)

func main() {
	text := "encrypt me"
	key := []byte("#89er@jdks$jmf_d")
	request := &restgoonch.Request{
		Text: text,
		Key:  key,
	}

	req, err := proto.Marshal(request)
	if err != nil {
		log.Fatalf("Unable to marshal request : %v", err)
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 200; i++ {
				resp, err := http.Post(fmt.Sprintf("http://%s:%d/service", host, port), "application/x-binary", bytes.NewReader(req))
				if err != nil {
					log.Fatalf("Unable to read from the server : %v", err)
				} /* else {
					log.Printf("Response from the server : %v", resp)
				}*/
				respBytes, err := ioutil.ReadAll(resp.Body)

				if err != nil {
					log.Fatalf("Unable to read bytes from request : %v", err)
				}
				respObj := &restgoonch.DecryptedText{}
				proto.Unmarshal(respBytes, respObj)
				log.Printf("Response from Goonch Server: %s", respObj.GetResult())
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"log"
	"strings"
)

var ES *elasticsearch.Client

const IndexName = "floors"

func InitSearch() {

	var r map[string]interface{}

	var err error
	ES, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{Config.ElasticsearchUrl},
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := ES.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err.Error())
	}

	// print Client and Server Info
	log.Printf("Client: %s\n", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	CheckIndex()
}

func CheckIndex() {

	rsp, err := ES.Indices.Get([]string{IndexName})
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	_ = rsp.Body.Close()
	if rsp.StatusCode == 404 {
		indexMapping := Map{
			"mappings": Map{
				"properties": Map{
					"content": Map{
						"type": "text",

						// https://www.elastic.co/guide/en/elasticsearch/reference/current/analyzer.html
						"analyzer": "ik_max_word",

						// https://www.elastic.co/guide/en/elasticsearch/reference/current/search-analyzer.html
						"search_analyzer": "ik_smart",

						// https://www.elastic.co/guide/en/elasticsearch/reference/current/similarity.html
						"similarity": "boolean",
					},
					"id": Map{
						// https://www.elastic.co/guide/en/elasticsearch/reference/current/number.html
						"type": "long",
					},
					"updated_at": Map{
						// https://www.elastic.co/guide/en/elasticsearch/reference/current/date.html
						"type": "date",
					},
				},
			},
		}
		buffer := bytes.NewBuffer(make([]byte, 0, 1024))
		err = json.NewEncoder(buffer).Encode(indexMapping)
		if err != nil {
			log.Fatal(err)
		}
		req := esapi.IndicesCreateRequest{
			Index: IndexName,
			Body:  buffer,
		}

		rsp, err = req.Do(context.Background(), ES)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}

		if rsp.IsError() {
			log.Fatalf("Error: %s", rsp.String())
		}
	}
}

var BulkBuffer *bytes.Buffer

// BulkInsert run in single goroutine only, used when dump floors
// see https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-bulk.html
func BulkInsert(floors Floors) error {
	if BulkBuffer == nil {
		BulkBuffer = bytes.NewBuffer(make([]byte, 0, 1024000)) // 100 KB buffer
	}

	if len(floors) == 0 {
		return nil
	}

	firstFloorID := floors[0].ID
	lastFloorID := floors[len(floors)-1].ID
	for _, floor := range floors {
		// meta: use index, it will insert or replace a document
		BulkBuffer.WriteString(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, floor.ID, "\n"))
		// data: should not contain \n, because \n is the delimiter of one action
		data, err := json.Marshal(floor)
		if err != nil {
			return fmt.Errorf("error failed to marshal floor: %s", err)
		}
		BulkBuffer.Write(data)
		BulkBuffer.WriteByte('\n') // the final line of data must end with a newline character \n
	}

	log.Printf("Preparing insert floor [%d, %d]\n", firstFloorID, lastFloorID)

	res, err := ES.Bulk(BulkBuffer, ES.Bulk.WithIndex(IndexName))
	if err != nil || res.IsError() {
		return fmt.Errorf("error indexing floor [%d, %d]: %s", firstFloorID, lastFloorID, err)
	}
	_ = res.Body.Close()
	log.Printf("index floor [%d, %d] success\n", firstFloorID, lastFloorID)

	BulkBuffer.Reset()
	return nil
}

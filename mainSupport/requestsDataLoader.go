package mainSupport

import (
	"endpointTester/httpRequest"
	"log"
	"os"
	"strings"
)

func LoadRequestsData(pathToRequestsDataConfig string) []*httpRequest.RequestDefinition {
	content, err := os.ReadFile(pathToRequestsDataConfig)
	if err != nil {
		log.Fatal(err)
	}
	docs := strings.Split(string(content), "---Doc---")
	requestsData := make([]*httpRequest.RequestDefinition, 0, len(docs))
	for index, doc := range docs {
		requestData, err := httpRequest.Parse(index, doc)
		if err != nil {
			log.Fatal("\n:::", index, "\n", doc, "\n", err)
		}
		if requestData != nil {
			requestsData = append(requestsData, requestData)
		}
	}
	return requestsData
}

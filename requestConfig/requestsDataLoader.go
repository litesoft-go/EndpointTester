package requestConfig

import (
	"log"
	"os"
	"strings"
)

func LoadRequestsData(pathToRequestsDataConfig string) []*RequestDefinition {
	content, err := os.ReadFile(pathToRequestsDataConfig)
	if err != nil {
		log.Fatal(err)
	}
	requestBlocks := strings.Split(":\n"+string(content), BlockStart)
	requestsDefs := make([]*RequestDefinition, 0, len(requestBlocks))
	startables := 0
	for index, block := range requestBlocks {
		requestDef, err := Parse(index, block)
		if err != nil {
			log.Fatal("\nblock[", index, "]\n", BlockStart, block, "\n", err)
		}
		if requestDef != nil {
			requestsDefs = append(requestsDefs, requestDef)
			if requestDef.IsSequenceStart() {
				startables++
			}
		}
	}
	if len(requestsDefs) == 0 {
		log.Fatal("\nNo block(s) converted to Request Definitions")
	}
	if startables == 0 {
		log.Fatal("\nAll block(s) are tagged with a Sequence continuation string (anything not empty or that doesn't spell 'sequence')")
	}
	return requestsDefs
}

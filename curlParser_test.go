package curlParser

import (
	"fmt"
	"testing"
)

func TestCurlParser_Parse(t *testing.T) {
	curlText := `curl --location --request GET 'https://localhost:8091/api?test=11' \
--header 'X-Header1: 1' \
--header 'X-Header2: 2' \
--header 'X-Header3: 3' \
--header 'Content-Type: application/json' \
--header 'X-Header4: 4'`

	curlParser := new(CurlParser)

	err, curlRequest := curlParser.Parse(curlText)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(curlRequest)
}

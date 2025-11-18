package pcurl

import (
	"fmt"
	"testing"

	"github.com/antlabs/pcurl/body"
)

func TestParseAndObj(t *testing.T) {
	// TODO 如果没有加右'接尾，报错
	all, err := ParseAndObj(`curl https://api.openai.com/v1/completions -H 'Content-Type: application/json' -H 'Authorization: Bearer YOUR_API_KEY' -d '{ "model": "text-davinci-003", "prompt": "Say this is a test", "max_tokens": 7, "temperature": 0 }'`)
	if err != nil {
		t.Fatalf("ParseAndObj failed: %v", err)
	}
	if all.Encode.Body != body.EncodeJSON {
		t.Fatalf("unexpected encode body, got=%v want=%v", all.Encode.Body, body.EncodeJSON)
	}
	if all.Method != "POST" {
		t.Fatalf("unexpected method, got=%q want=%q", all.Method, "POST")
	}
	fmt.Printf("%#v\n", all)
}

func TestParseAndJSON(t *testing.T) {
	// TODO 如果没有加右'接尾，报错
	all, err := ParseAndJSON(`curl https://api.openai.com/v1/completions -H 'Content-Type: application/json' -H 'Authorization: Bearer YOUR_API_KEY' -d '{ "model": "text-davinci-003", "prompt": "Say this is a test", "max_tokens": 7, "temperature": 0 }'`)
	if err != nil {
		t.Fatalf("ParseAndJSON failed: %v", err)
	}
	fmt.Printf("%s\n", all)
}

type testCaseObj struct {
	curl string
	url  string
	body string
}

func TestPaserAndObj(t *testing.T) {

	tab := []testCaseObj{
		{
			curl: `curl -X POST -d '{"a":"b"}' www.qq.com/test`,
			url:  `www.qq.com/test`,
		},
		{
			curl: `curl -X POST -i 'http://{{.Host}}/{{.OrgName}}/{{.AppName}}/messages/users' -H 'Content-Type: application/json' -H 'Accept: application/json' -H 'Authorization: Bearer <YourAppToken>' -d '{"from": "user1","to": ["user2"],"type": "txt","body": {"msg": "testmessages"}}'`,
			url:  `http://{{.Host}}/{{.OrgName}}/{{.AppName}}/messages/users`,
		},
		{
			curl: `curl -X DELETE -H 'Accept: application/json' -H 'Authorization: Bearer <YourAppToken> ' https://{{.Host}}/{{.OrgName}}/{{.AppName}}/chatgroups/{{.GroupID}}`,
			url:  `https://{{.Host}}/{{.OrgName}}/{{.AppName}}/chatgroups/{{.GroupID}}`,
		},
	}

	for _, tc := range tab {

		all, err := ParseAndObj(tc.curl)
		if err != nil {
			t.Fatalf("ParseAndObj failed for curl=%q: %v", tc.curl, err)
		}
		if all.URL != tc.url {
			t.Fatalf("unexpected URL, got=%q want=%q", all.URL, tc.url)
		}
	}
}

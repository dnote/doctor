package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

/** snapshots **/

type dnoteV0_4_0V0_4_4 map[string]bookV0_4_0V0_4_4

type bookV0_4_0V0_4_4 struct {
	Name  string             `json:"name"`
	Notes []noteV0_4_0V0_4_4 `json:"notes"`
}

type noteV0_4_0V0_4_4 struct {
	UUID     string `json:"uuid"`
	Content  string `json:"content"`
	AddedOn  int64  `json:"added_on"`
	EditedOn int64  `json:"edited_on"`
	Public   bool   `json:"public"`
}

func readJSONDnote(ctx Ctx) (json.RawMessage, error) {
	var books json.RawMessage

	notePath := fmt.Sprintf("%s/.dnote/dnote", ctx.homeDirPath)
	b, err := ioutil.ReadFile(notePath)
	if err != nil {
		return nil, errors.Wrap(err, "reading note content")
	}

	err = json.Unmarshal(b, &books)
	if err != nil {
		return books, errors.Wrap(err, "unmarshalling note content")
	}

	return books, nil
}

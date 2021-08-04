package pretend

import (
	"encoding/json"
	"fmt"

	"github.com/hypha-dao/document-graph/docgraph"
)

func DefaultSettings() ([]docgraph.ContentItem, error) {

	var settings []docgraph.ContentItem
	err := json.Unmarshal([]byte(rawSettings), &settings)
	if err != nil {
		return settings, fmt.Errorf("cannot unmarshal error: %v ", err)
	}
	return settings, nil
}

const rawSettings = `[
	  {
		"label": "husd_token_contract",
		"value": [
		  "name",
		  "husd.hypha"
		]
	  },{
		"label": "hypha_token_contract",
		"value": [
		  "name",
		  "token.hypha"
		]
	  },{
		"label": "last_ballot_id",
		"value": [
		  "name",
		  "hypha1....1tg"
		]
	  },{
		"label": "publisher_contract",
		"value": [
		  "name",
		  "publsh.hypha"
		]
	  },{
		"label": "seeds_escrow_contract",
		"value": [
		  "name",
		  "escrow.seeds"
		]
	  },{
		"label": "seeds_token_contract",
		"value": [
		  "name",
		  "token.seeds"
		]
	  },{
		"label": "telos_decide_contract",
		"value": [
		  "name",
		  "trailservice"
		]
	  },{
		"label": "hvoice_token_contract",
		"value": [
		  "name",
		  "voice.hypha"
		]
	  },{
		"label": "treasury_contract",
		"value": [
		  "name",
		  "bank.hypha"
		]
	  },{
		"label": "client_version",
		"value": [
		  "string",
		  "1.0.13 0c81dde6"
		]
	  },{
		"label": "contract_version",
		"value": [
		  "string",
		  "v0.2.0 618f051"
		]
	  },{
		"label": "hypha_deferral_factor_x100",
		"value": [
		  "int64",
		  13
		]
	  },{
		"label": "last_sender_id",
		"value": [
		  "int64",
		  527
		]
	  },{
		"label": "paused",
		"value": [
		  "int64",
		  0
		]
	  },{
		"label": "seeds_deferral_factor_x100",
		"value": [
		  "int64",
		  100
		]
	  }
	]`

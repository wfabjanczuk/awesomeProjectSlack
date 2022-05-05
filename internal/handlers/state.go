package handlers

import (
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/requests"
)

const PublicChannel = "public"

var requestQueue = make(chan *requests.WSRequest)
var channels = make(map[string][]*WSConnection)

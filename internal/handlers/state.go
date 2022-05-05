package handlers

import (
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/models"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/requests"
)

var publicChannel = &models.Channel{
	Name: "public",
}
var requestQueue = make(chan *requests.WSRequest)
var channels = make(map[models.Channel][]WSConnection)

package handlers

import (
	"fmt"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/connections"
)

func broadcastToChannel(channelName, message string) {
	for clientConnection := range channels[channelName] {
		sendSuccessMessage(clientConnection, message)
	}
}

func createChannel(channelName string, clientConnection *connections.ClientConnection) {
	if _, ok := channels[channelName]; ok {
		errorMessage := fmt.Sprintf("Channel with name: \"%s\" already exists!", channelName)
		sendErrorMessage(clientConnection, errorMessage)

		return
	}

	channels[channelName] = NewClientConnectionSet()

	successMessage := fmt.Sprintf("Channel with name: \"%s\" successfully created.", channelName)
	sendSuccessMessage(clientConnection, successMessage)
}

func enterChannel(channelName string, clientConnection *connections.ClientConnection) {
	if _, ok := channels[channelName]; !ok {
		errorMessage := fmt.Sprintf("Channel with name: \"%s\" does not exist!", channelName)
		sendErrorMessage(clientConnection, errorMessage)

		return
	}

	lastChannel := clientConnection.GetChannel()
	if channelName == lastChannel {
		successMessage := fmt.Sprintf("You already are on channel with name: \"%s\"", channelName)
		sendSuccessMessage(clientConnection, successMessage)

		return
	}

	channels[clientConnection.GetChannel()].Delete(clientConnection)
	clientConnection.SetChannel(channelName)
	channels[channelName].Add(clientConnection)

	successMessage := fmt.Sprintf("Successfully left channel with name: \"%s\" and entered channel \"%s\"", lastChannel, channelName)
	sendSuccessMessage(clientConnection, successMessage)
}

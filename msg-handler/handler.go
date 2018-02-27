/*
 * Copyright (c) 2018 Juniper Networks, Inc. All rights reserved.
 *
 * file:    handler.go
 * details: Initializes all the message handlers
 *
 */
package msghandler

import (
	opts "github.com/Juniper/ipfix-translator/options"
	"sync"
)

type Handler struct {
	MH     MsgHandler
	MHChan chan []byte
}

type MsgHandler interface {
	setup(string) error
	handleMessages(chan []byte)
}

func NewMsgHandler(handlerName string) *Handler {
	var msgHandlerRegistered = map[string]MsgHandler{
		"data-manager": new(DataManager),
	}
	return &Handler{
		MH: msgHandlerRegistered[handlerName],
	}
}

func (h Handler) Run() error {
	var (
		wg  sync.WaitGroup
		err error
	)
	err = h.MH.setup(opts.MHConfigFile)
	if err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		h.MH.handleMessages(h.MHChan)
	}()

	wg.Wait()

	return nil
}

package rpc

import (
	"goproject/rpc/codec"
	"io"
	"log"
	"reflect"
)

var invalidRequest = struct{}{}

type request struct {
	h            *codec.Header
	argv, replyv reflect.Value // 请求的Argv和reply
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Printf("rpc readRequestHeader err: %v\n", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	var r = &request{h: h}
	r.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(r.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return r, nil
}

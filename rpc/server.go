package rpc

import (
	"encoding/json"
	"fmt"
	"goproject/rpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int
	CodecType   codec.Type
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go server.ServeConn(conn)
	}
}

func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Printf("rpc ServeConn err: %v\n", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc ServeConn magic num err: %v\n", opt.MagicNumber)
		return
	}
	gen, ok := codec.NewCodecFuncMap[opt.CodecType]
	if !ok {
		log.Printf("rcp ServeConn unsupported codec: %v\n", opt.CodecType)
		return
	}
	server.serveCodec(gen(conn))
}

func (server *Server) serveCodec(cc codec.Codec) {
	wg := new(sync.WaitGroup)
	sending := new(sync.Mutex)
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Printf("rpc sendResponse err: %v\n", err)
	}
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("rpc handleRequest： %v\t%v\n", req.h, req.argv)
	req.replyv = reflect.ValueOf(fmt.Sprintf("rpc resp %d", req.h.Seq))
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}

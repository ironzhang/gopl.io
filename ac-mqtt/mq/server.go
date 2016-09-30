package mq

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/surgemq/message"

	"ac-common-go/glog"
	"ac-common-go/net/context"
	"ac-mqtt/mq/session"
	_ "ac-mqtt/mq/session/mem-session-map"
)

const (
	minKeepAlive = 30
)

type Authenticator interface {
	Authenticate(user, password []byte) error
}

type Server struct {
	net  string
	addr string
	ln   net.Listener
	smap session.Map
	auth Authenticator

	rbsize         int
	wbsize         int
	connectTimeout time.Duration
}

func NewServer(net, addr string) *Server {
	return new(Server).Init(net, addr)
}

func (s *Server) Init(net, addr string) *Server {
	s.net = net
	s.addr = addr
	s.smap, _ = session.NewMap("mem-session-map", "")
	s.rbsize = 1024
	s.wbsize = 1024
	s.connectTimeout = time.Microsecond
	return s
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen(s.net, s.addr)
	if err != nil {
		glog.Errorf("Listen(%q, %q): %v", s.net, s.addr, err)
		return err
	}
	s.ln = ln

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxid := 0
	glog.V(3).Infof("server(%s/%s) start serve", s.net, s.addr)
	for {
		nc, err := ln.Accept()
		if err != nil {
			glog.V(1).Infof("accept: %v", err)
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				glog.V(2).Infof("accept occur temporary error: %v", err)
				continue
			}
			break
		}

		maxid++
		c := newconn(maxid, nc)
		go s.handleConn(ctx, c)
	}
	glog.V(3).Infof("server(%s/%s) stop serve", s.net, s.addr)

	return nil
}

func (s *Server) Close() error {
	s.ln.Close()
	return nil
}

func (s *Server) handleConn(ctx context.Context, c *conn) {
	var err error
	glog.V(4).Infof("conn(%s) handle start", c.String())
	defer func() {
		glog.V(4).Infof("conn(%s) handle stop: %v", c.String(), err)
	}()

	rch := make(chan message.Message, s.rbsize)
	wch := make(chan message.Message, s.wbsize)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			close(rch)
			close(wch)
			c.Close()
		}
	}()

	// 处理第一包连接消息
	if err = s.dealConnectMessage(c); err != nil {
		return
	}

	// 启动消息处理goroutine
	go s.processor(c, rch)

	// 启动消息写入goroutine
	go s.writer(c, wch)

	// 读取消息
	err = s.reader(c, rch)
}

func (s *Server) reader(c *conn, msgc chan<- message.Message) error {
	var err error
	var msg message.Message

	for {
		msg, err = readMessage(c)
		if err != nil {
			if err == io.EOF {
				break
			}

			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				glog.V(1).Infof("conn(%s) read message occur temporary error: %v", c.String(), err)
				continue
			}

			break
		}
		msgc <- msg
	}
	glog.V(4).Infof("conn(%s) reader stop: %v", c.String(), err)

	return err
}

func (s *Server) writer(c *conn, msgc <-chan message.Message) {
	var err error
	for msg := range msgc {
		if err = writeMessage(c, msg); err != nil {
			glog.V(1).Infof("write message: %v", err)
			continue
		}
	}
	glog.V(4).Infof("conn(%s) writer stop", c.String())
}

func (s *Server) processor(c *conn, msgc <-chan message.Message) {
	for msg := range msgc {
		glog.V(4).Infof("conn(%s) process msg. Name:%s, PacketId:%d", c.String(), msg.Name(), msg.PacketId())
		if err := s.processMessage(c, msg); err != nil {
			glog.V(1).Infof("conn(%s) process message(%s): %v", c.String(), msg.Name(), err)
		}
	}
	glog.V(4).Infof("conn(%s) processor stop", c.String())
}

func (s *Server) dealConnectMessage(c *conn) error {
	resp := message.NewConnackMessage()

	// 读取连接消息
	if s.connectTimeout > 0 {
		c.SetReadDeadline(time.Now().Add(s.connectTimeout))
	}
	req, err := readConnectMessage(c)
	if err != nil {
		glog.V(1).Infof("conn(%s) read connect message: %v", c.String(), err)
		if cerr, ok := err.(message.ConnackCode); ok {
			resp.SetReturnCode(cerr)
			resp.SetSessionPresent(false)
			writeMessage(c, resp)
		}
		return err
	}
	if s.connectTimeout > 0 {
		c.SetReadDeadline(time.Time{})
	}

	// 认证客户端
	if s.auth != nil {
		if err = s.auth.Authenticate(req.Username(), req.Password()); err != nil {
			glog.V(3).Infof("conn(%s) Authenticate(%q, %q): %v", c.String(), req.Username(), req.Password(), err)
			resp.SetReturnCode(message.ErrBadUsernameOrPassword)
			resp.SetSessionPresent(false)
			writeMessage(c, resp)
			return err
		}
	}

	// 关联session
	if err = s.attachSession(c, req, resp); err != nil {
		glog.V(1).Infof("conn(%s) attach session: %v", c.String(), err)
		if cerr, ok := err.(message.ConnackCode); ok {
			resp.SetReturnCode(cerr)
			resp.SetSessionPresent(false)
			writeMessage(c, resp)
		}
		return err
	}

	// 发送连接回执
	resp.SetReturnCode(message.ConnectionAccepted)
	if err = writeMessage(c, resp); err != nil {
		glog.V(1).Infof("conn(%s) write connack message: %v", c.String(), err)
		return err
	}

	return nil
}

func (s *Server) attachSession(c *conn, req *message.ConnectMessage, resp *message.ConnackMessage) error {
	id := string(req.ClientId())
	if id == "" {
		id = fmt.Sprintf("c%d", c.id)
		resp.SetSessionPresent(false)
		return s.newSession(id, req)
	}

	if req.CleanSession() {
		resp.SetSessionPresent(false)
		return s.newSession(id, req)
	}

	if _, ok := s.smap.Get(id); !ok {
		resp.SetSessionPresent(false)
		return s.newSession(id, req)
	}

	resp.SetSessionPresent(true)
	return nil
}

func (s *Server) newSession(id string, req *message.ConnectMessage) error {
	glog.V(4).Infof("new session: %s", id)
	sess := session.New(id)
	if req.WillFlag() {
		sess.SetWillMessage(req.WillQos(), req.WillTopic(), req.WillMessage(), req.WillRetain())
	}
	return s.smap.Set(sess)
}

func (s *Server) processMessage(c *conn, m message.Message) error {
	switch msg := m.(type) {
	case *message.PublishMessage:
		// 发布消息
		return s.processPublish(c, msg)

	case *message.PubrelMessage:
		// QoS2消息释放
		return s.processPubrel(c, msg)

	case *message.SubscribeMessage:
		// 订阅主题
		return s.processSubscribe(c, msg)

	case *message.UnsubscribeMessage:
		// 取消订阅
		return s.processUnsubscribe(c, msg)

	case *message.PingreqMessage:
		// PING请求
		return s.processPingreq(c, msg)

	case *message.DisconnectMessage:
		// 断开连接
		c.Close()
		return nil
	}
	return nil
}

func (s *Server) processPublish(c *conn, msg *message.PublishMessage) error {
	switch msg.QoS() {
	// 不需要响应
	case message.QosAtMostOnce:
		//glog.V(4).Infof("conn(%s) publish message.QosAtMostOnce", c.String())
		return nil

	// 响应PUBACK
	case message.QosAtLeastOnce:
		//glog.V(4).Infof("conn(%s) publish message.QosAtLeastOnce", c.String())
		resp := message.NewPubackMessage()
		resp.SetPacketId(msg.PacketId())
		if err := writeMessage(c, resp); err != nil {
			return err
		}
		return nil

	// 响应PUBREC
	case message.QosExactlyOnce:
		//glog.V(4).Infof("conn(%s) publish message.QosExactlyOnce", c.String())
		resp := message.NewPubrecMessage()
		resp.SetPacketId(msg.PacketId())
		if err := writeMessage(c, resp); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (s *Server) processPubrel(c *conn, msg *message.PubrelMessage) error {
	resp := message.NewPubcompMessage()
	resp.SetPacketId(msg.PacketId())
	if err := writeMessage(c, resp); err != nil {
		return err
	}
	return nil
}

func (s *Server) processSubscribe(c *conn, msg *message.SubscribeMessage) error {
	resp := message.NewSubackMessage()
	resp.SetPacketId(msg.PacketId())

	topics := msg.Topics()
	for range topics {
		resp.AddReturnCode(message.QosAtMostOnce)
	}

	if err := writeMessage(c, resp); err != nil {
		return err
	}
	return nil
}

func (s *Server) processUnsubscribe(c *conn, msg *message.UnsubscribeMessage) error {
	resp := message.NewUnsubackMessage()
	resp.SetPacketId(msg.PacketId())
	if err := writeMessage(c, resp); err != nil {
		return err
	}
	return nil
}

func (s *Server) processPingreq(c *conn, msg *message.PingreqMessage) error {
	resp := message.NewPingrespMessage()
	if err := writeMessage(c, resp); err != nil {
		return err
	}
	return nil
}

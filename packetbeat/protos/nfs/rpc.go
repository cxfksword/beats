// Package rpc provides support for parsing RPC messages and reporting the
// results. This package supports the RPC v2 protocol as defined by RFC 5531
// (RFC 1831).

package nfs

import (
	"encoding/binary"
	"time"

	"github.com/cxfksword/beats/libbeat/common"
	"github.com/cxfksword/beats/libbeat/logp"

	"github.com/cxfksword/beats/packetbeat/protos"
	"github.com/cxfksword/beats/packetbeat/protos/tcp"
	"github.com/cxfksword/beats/packetbeat/publish"
)

var debugf = logp.MakeDebug("rpc")

const (
	RPC_LAST_FRAG = 0x80000000
	RPC_SIZE_MASK = 0x7fffffff
)

type RpcStream struct {
	tcpTuple *common.TcpTuple
	rawData  []byte
}

type rpcConnectionData struct {
	Streams [2]*RpcStream
}

type Rpc struct {
	// Configuration data.
	Ports []int

	transactionTimeout time.Duration

	results publish.Transactions // Channel where results are pushed.
}

func init() {
	protos.Register("nfs", New)
}

func New(
	testMode bool,
	results publish.Transactions,
	cfg *common.Config,
) (protos.Plugin, error) {
	p := &Rpc{}
	config := defaultConfig
	if !testMode {
		if err := cfg.Unpack(&config); err != nil {
			logp.Warn("failed to read config")
			return nil, err
		}
	}

	if err := p.init(results, &config); err != nil {
		logp.Warn("failed to init")
		return nil, err
	}
	return p, nil
}

func (rpc *Rpc) init(results publish.Transactions, config *rpcConfig) error {
	rpc.setFromConfig(config)
	rpc.results = results

	return nil
}

func (rpc *Rpc) setFromConfig(config *rpcConfig) error {
	rpc.Ports = config.Ports
	rpc.transactionTimeout = time.Duration(config.TransactionTimeout) * time.Second
	return nil
}

func (rpc *Rpc) GetPorts() []int {
	return rpc.Ports
}

// Called when TCP payload data is available for parsing.
func (rpc *Rpc) Parse(
	pkt *protos.Packet,
	tcptuple *common.TcpTuple,
	dir uint8,
	private protos.ProtocolData,
) protos.ProtocolData {

	defer logp.Recover("ParseRPC exception")

	conn := ensureRpcConnection(private)

	conn = rpc.handleRpcFragment(conn, pkt, tcptuple, dir)
	if conn == nil {
		return nil
	}
	return conn
}

// Called when the FIN flag is seen in the TCP stream.
func (rpc *Rpc) ReceivedFin(tcptuple *common.TcpTuple, dir uint8,
	private protos.ProtocolData) protos.ProtocolData {

	defer logp.Recover("ReceivedFinRpc exception")

	// forced by TCP interface

	// TODO
	return private
}

// Called when a packets are missing from the tcp
// stream.
func (rpc *Rpc) GapInStream(tcptuple *common.TcpTuple, dir uint8,
	nbytes int, private protos.ProtocolData) (priv protos.ProtocolData, drop bool) {

	defer logp.Recover("GapInRpcStream exception")

	// forced by TCP interface

	// TODO
	return private, false
}

// ConnectionTimeout returns the per stream connection timeout.
// Return <=0 to set default tcp module transaction timeout.
func (rpc *Rpc) ConnectionTimeout() time.Duration {
	// forced by TCP interface
	return rpc.transactionTimeout
}

func ensureRpcConnection(private protos.ProtocolData) *rpcConnectionData {
	conn := getRpcConnection(private)
	if conn == nil {
		conn = &rpcConnectionData{}
	}
	return conn
}

func getRpcConnection(private protos.ProtocolData) *rpcConnectionData {
	if private == nil {
		return nil
	}

	priv, ok := private.(*rpcConnectionData)
	if !ok {
		logp.Warn("rpc connection data type error")
		return nil
	}
	if priv == nil {
		logp.Warn("Unexpected: rpc connection data not set")
		return nil
	}

	return priv
}

// Parse function is used to process TCP payloads.
func (rpc *Rpc) handleRpcFragment(
	conn *rpcConnectionData,
	pkt *protos.Packet,
	tcptuple *common.TcpTuple,
	dir uint8,
) *rpcConnectionData {

	st := conn.Streams[dir]
	if st == nil {
		st = newStream(pkt, tcptuple)
		conn.Streams[dir] = st
	} else {
		// concatenate bytes
		st.rawData = append(st.rawData, pkt.Payload...)
		if len(st.rawData) > tcp.TCP_MAX_DATA_IN_STREAM {
			debugf("Stream data too large, dropping TCP stream")
			conn.Streams[dir] = nil
			return conn
		}
	}

	for len(st.rawData) > 0 {

		if len(st.rawData) < 4 {
			debugf("Wainting for more data")
			break
		}

		marker := uint32(binary.BigEndian.Uint32(st.rawData[0:4]))
		size := int(marker & RPC_SIZE_MASK)
		islast := (marker & RPC_LAST_FRAG) != 0

		if len(st.rawData)-4 < size {
			debugf("Wainting for more data")
			break
		}

		if !islast {
			logp.Warn("multifragment rpc message")
			break
		}

		xdr := Xdr{data: st.rawData[4 : 4+size], offset: 0}
		msg := &RpcMessage{ts: pkt.Ts, xdr: xdr}

		src := common.Endpoint{
			Ip:   tcptuple.Src_ip.String(),
			Port: tcptuple.Src_port,
		}
		dst := common.Endpoint{
			Ip:   tcptuple.Dst_ip.String(),
			Port: tcptuple.Dst_port,
		}

		event := common.MapStr{}
		event["@timestamp"] = common.Time(pkt.Ts)
		event["status"] = common.OK_STATUS // all packes are OK for now
		event["src"] = &src
		event["dst"] = &dst

		msg.fillEvent(event, rpc.results, size)

		st.rawData = st.rawData[4+size:]
	}

	return conn
}

func newStream(pkt *protos.Packet, tcptuple *common.TcpTuple) *RpcStream {
	return &RpcStream{
		tcpTuple: tcptuple,
	}
}

//

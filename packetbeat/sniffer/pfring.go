// +build linux,havepfring

package sniffer

import (
	"fmt"

	"github.com/cxfksword/gopacket"
	"github.com/cxfksword/gopacket/pfring"
)

type PfringHandle struct {
	Ring *pfring.Ring
}

func NewPfringHandle(device string, snaplen int, promisc bool) (*PfringHandle, error) {

	var h PfringHandle
	var err error

	if device == "any" {
		return nil, fmt.Errorf("Pfring sniffing doesn't support 'any' as interface")
	}

	var flags pfring.Flag

	if promisc {
		flags = pfring.FlagPromisc
	}

	h.Ring, err = pfring.NewRing(device, uint32(snaplen), flags)

	return &h, err
}

func (h *PfringHandle) ReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error) {
	return h.Ring.ReadPacketData()
}

func (h *PfringHandle) SetBPFFilter(expr string) (_ error) {
	return h.Ring.SetBPFFilter(expr)
}

func (h *PfringHandle) Enable() (_ error) {
	return h.Ring.Enable()
}

func (h *PfringHandle) Close() {
	h.Ring.Close()
}

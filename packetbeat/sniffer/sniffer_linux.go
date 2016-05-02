// +build linux

package sniffer

import (
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/cxfksword/beats/libbeat/logp"

	"github.com/cxfksword/beats/packetbeat/config"

	"github.com/tsg/gopacket"
	"github.com/tsg/gopacket/layers"
)

type SnifferSetup struct {
	afpacketHandle *AfpacketHandle
	config         *config.InterfacesConfig
	isAlive        bool

	// bpf filter
	filter string

	// Decoder    *decoder.DecoderStruct
	worker     Worker
	DataSource gopacket.PacketDataSource
}

type Worker interface {
	OnPacket(data []byte, ci *gopacket.CaptureInfo)
}

type WorkerFactory func(layers.LinkType) (Worker, string, error)

// Computes the block_size and the num_blocks in such a way that the
// allocated mmap buffer is close to but smaller than target_size_mb.
// The restriction is that the block_size must be divisible by both the
// frame size and page size.
func afpacketComputeSize(target_size_mb int, snaplen int, page_size int) (
	frame_size int, block_size int, num_blocks int, err error) {

	if snaplen < page_size {
		frame_size = page_size / (page_size / snaplen)
	} else {
		frame_size = (snaplen/page_size + 1) * page_size
	}

	// 128 is the default from the gopacket library so just use that
	block_size = frame_size * 128
	num_blocks = (target_size_mb * 1024 * 1024) / block_size

	if num_blocks == 0 {
		return 0, 0, 0, fmt.Errorf("Buffer size too small")
	}

	return frame_size, block_size, num_blocks, nil
}

func deviceNameFromIndex(index int, devices []string) (string, error) {
	if index >= len(devices) {
		return "", fmt.Errorf("Looking for device index %d, but there are only %d devices",
			index, len(devices))
	}

	return devices[index], nil
}

// ListDevicesNames returns the list of adapters available for sniffing on
// this computer. If the withDescription parameter is set to true, a human
// readable version of the adapter name is added.
func ListDeviceNames(withDescription bool) ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return []string{}, err
	}

	ret := []string{}
	for _, iface := range ifaces {
		if strings.Contains(iface.Name, "eth") {
			if withDescription {
				desc := "No description available"
				ret = append(ret, fmt.Sprintf("%s (%s)", iface.Name, desc))
			} else {
				ret = append(ret, iface.Name)
			}
		}
	}

	if withDescription {
		ret = append(ret, fmt.Sprintf("%s (%s)", "any", "Pseudo-device that captures on all interfaces"))
		ret = append(ret, fmt.Sprintf("%s (%s)", "lo", "No description available"))
	} else {
		ret = append(ret, "any")
		ret = append(ret, "lo")
	}

	return ret, nil
}

func (sniffer *SnifferSetup) setFromConfig(config *config.InterfacesConfig) error {
	sniffer.config = config

	if len(sniffer.config.File) > 0 {
		logp.Debug("sniffer", "Reading from file: %s", sniffer.config.File)
		// we read file with the pcap provider
		sniffer.config.Type = "pcap"
	}

	// set defaults
	if len(sniffer.config.Device) == 0 {
		sniffer.config.Device = "any"
	}

	if index, err := strconv.Atoi(sniffer.config.Device); err == nil { // Device is numeric
		devices, err := ListDeviceNames(false)
		if err != nil {
			return fmt.Errorf("Error getting devices list: %v", err)
		}
		sniffer.config.Device, err = deviceNameFromIndex(index, devices)
		if err != nil {
			return fmt.Errorf("Couldn't understand device index %d: %v", index, err)
		}
		logp.Info("Resolved device index %d to device: %s", index, sniffer.config.Device)
	}

	if sniffer.config.Snaplen == 0 {
		sniffer.config.Snaplen = 65535
	}

	if sniffer.config.Type == "autodetect" || sniffer.config.Type == "" {
		if runtime.GOOS == "linux" {
			sniffer.config.Type = "af_packet"
		} else {
			sniffer.config.Type = "pcap"
		}
	}

	logp.Info("Sniffer type: %s device: %s BPF filter: '%s'", sniffer.config.Type, sniffer.config.Device, sniffer.filter)

	switch sniffer.config.Type {

	case "af_packet":
		if sniffer.config.Buffer_size_mb == 0 {
			sniffer.config.Buffer_size_mb = 24
		}

		frame_size, block_size, num_blocks, err := afpacketComputeSize(
			sniffer.config.Buffer_size_mb,
			sniffer.config.Snaplen,
			os.Getpagesize())
		if err != nil {
			return err
		}

		sniffer.afpacketHandle, err = NewAfpacketHandle(
			sniffer.config.Device,
			frame_size,
			block_size,
			num_blocks,
			500*time.Millisecond)
		if err != nil {
			return err
		}

		err = sniffer.afpacketHandle.SetBPFFilter(sniffer.filter)
		if err != nil {
			return fmt.Errorf("SetBPFFilter failed: %s", err)
		}

		sniffer.DataSource = gopacket.PacketDataSource(sniffer.afpacketHandle)

	default:
		return fmt.Errorf("Unknown sniffer type: %s", sniffer.config.Type)
	}

	return nil
}

func (sniffer *SnifferSetup) Reopen() error {
	if sniffer.config.Type != "pcap" || sniffer.config.File == "" {
		return fmt.Errorf("Reopen is only possible for files")
	}

	return nil
}

func (sniffer *SnifferSetup) Datalink() layers.LinkType {
	return layers.LinkTypeEthernet
}

func (sniffer *SnifferSetup) Init(test_mode bool, factory WorkerFactory, filter string) error {
	var err error

	sniffer.filter = filter
	if !test_mode {
		err = sniffer.setFromConfig(&config.ConfigSingleton.Interfaces)
		if err != nil {
			return fmt.Errorf("Error creating sniffer: %v", err)
		}
	}

	sniffer.worker, sniffer.filter, err = factory(sniffer.Datalink())
	if err != nil {
		return fmt.Errorf("Error creating decoder: %v", err)
	}

	if sniffer.config.Dumpfile != "" {
		return fmt.Errorf("Linux not support Dumpfile")
	}

	sniffer.isAlive = true

	return nil
}

func (sniffer *SnifferSetup) Run() error {
	counter := 0
	loopCount := 1
	var lastPktTime *time.Time = nil
	var ret_error error
	for sniffer.isAlive {
		if sniffer.config.OneAtATime {
			fmt.Println("Press enter to read packet")
			fmt.Scanln()
		}

		data, ci, err := sniffer.DataSource.ReadPacketData()

		if err == syscall.EINTR {
			logp.Debug("sniffer", "Interrupted")
			continue
		}

		if err == io.EOF {
			logp.Debug("sniffer", "End of file")
			loopCount += 1
			if sniffer.config.Loop > 0 && loopCount > sniffer.config.Loop {
				// give a bit of time to the publish goroutine
				// to flush
				time.Sleep(300 * time.Millisecond)
				sniffer.isAlive = false
				continue
			}

			logp.Debug("sniffer", "Reopening the file")
			err = sniffer.Reopen()
			if err != nil {
				ret_error = fmt.Errorf("Error reopening file: %s", err)
				sniffer.isAlive = false
				continue
			}
			lastPktTime = nil
			continue
		}

		if err != nil {
			ret_error = fmt.Errorf("Sniffing error: %s", err)
			sniffer.isAlive = false
			continue
		}

		if len(data) == 0 {
			// Empty packet, probably timeout from afpacket
			continue
		}

		if sniffer.config.File != "" {
			if lastPktTime != nil && !sniffer.config.TopSpeed {
				sleep := ci.Timestamp.Sub(*lastPktTime)
				if sleep > 0 {
					time.Sleep(sleep)
				} else {
					logp.Warn("Time in pcap went backwards: %d", sleep)
				}
			}
			_lastPktTime := ci.Timestamp
			lastPktTime = &_lastPktTime
			if !sniffer.config.TopSpeed {
				ci.Timestamp = time.Now() // overwrite what we get from the pcap
			}
		}
		counter++

		logp.Debug("sniffer", "Packet number: %d", counter)

		sniffer.worker.OnPacket(data, &ci)
	}

	logp.Info("Input finish. Processed %d packets. Have a nice day!", counter)

	return ret_error
}

func (sniffer *SnifferSetup) Close() error {
	switch sniffer.config.Type {
	case "af_packet":
		sniffer.afpacketHandle.Close()
	}
	return nil
}

func (sniffer *SnifferSetup) Stop() error {
	sniffer.isAlive = false
	return nil
}

func (sniffer *SnifferSetup) IsAlive() bool {
	return sniffer.isAlive
}

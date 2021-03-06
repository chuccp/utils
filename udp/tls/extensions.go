package tls

import (
	"github.com/chuccp/utils/udp/config"
	"github.com/chuccp/utils/udp/util"
)

const (
	KeyShareType uint16 = 51
	TransportType uint16 = 57
)

type Extensions struct {
	ExtensionsMap map[uint16]*Extension
}

func (es *Extensions) SetKeyShare(KeyExchanges []byte) {
	ks := NewKeyShare(x25519,KeyExchanges)
	kWR:=  util.NewWriteBuffer()
	ks.Write(kWR)
	ex := NewExtension(KeyShareType, kWR.Bytes())
	es.addExtensions(ex)
}
func (es *Extensions)GetKeyShare(clientKeyShare *ClientKeyShare) error {
	extension := es.ExtensionsMap[KeyShareType]
	if extension==nil{
		return util.FormatError
	}
	err := UnPacketClientKeyShare(extension.Data, clientKeyShare)
	if err != nil {
		return err
	}
	return nil
}


func (es *Extensions) SetTransportParameters(sendConfig *config.SendConfig) {
	ntp:=NewTransportParameters()
	ntp.SetValue(MaxUdpPayloadSizeType,sendConfig.MaxUdpPayloadSize)
	ntp.SetValue(InitialMaxStreamDataBidiRemoteType,sendConfig.InitialMaxStreamDataBidiRemote)
	ntp.SetValue(InitialMaxStreamDataBidiLocalType,sendConfig.InitialMaxStreamDataBidiLocal)
	ntp.SetValue(InitialMaxDataType,sendConfig.InitialMaxData)
	ntp.SetValue(MaxIdleTimeout,sendConfig.MaxIdleTimeout)
	ntp.SetValue(InitialMaxStreamsBidi,sendConfig.InitialMaxStreamsBidi)
	ntp.SetValue(InitialMaxStreamsUni,sendConfig.InitialMaxStreamsUni)
	ntpWR:=  util.NewWriteBuffer()
	ntp.Write(ntpWR)
	ex := NewExtension(TransportType, ntpWR.Bytes())
	es.addExtensions(ex)
}


func (es *Extensions) addExtensions(ex *Extension) {
	es.ExtensionsMap[ex.Type] = ex
}

func NewExtensions() *Extensions {
	return &Extensions{ExtensionsMap: make(map[uint16]*Extension)}
}

func (es *Extensions) Write(write *util.WriteBuffer) {
	for _, extension := range es.ExtensionsMap {
		write.WriteUint16(extension.Type)
		write.WriteUint16LengthBuff(func(wr *util.WriteBuffer) {
			wr.WriteBytes(extension.Data)
		})
	}
}
func (es *Extensions) Read(read *util.ReadBuffer) error{

	for{
		ExtensionType, err := read.ReadUint16Length()
		if err != nil {
			return err
		}
		_, data, err := read.ReadU16LengthBytes()
		if err != nil {
			return err
		}
		ex:=NewExtension(ExtensionType,data)
		es.addExtensions(ex)
		if read.Buffered()==0{
			break
		}
	}
	return nil
}
func  ReadExtensions(read *util.ReadBuffer) (*Extensions,error){
	 extensions:=NewExtensions()
	return extensions,extensions.Read(read)
}
type Extension struct {
	Type   uint16
	Data []byte
}

func NewExtension(Type uint16, Data []byte) *Extension {
	return &Extension{Type: Type, Data: Data}
}
func (e *Extension) Write(write *util.WriteBuffer) {
	write.WriteUint16(e.Type)
	write.WriteUint16LengthBuff(func(write *util.WriteBuffer) {
		write.WriteBytes(e.Data)
	})
}

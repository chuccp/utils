package wire

import (
	"github.com/chuccp/utils/udp/util"
)

func UnPacketCryptoFrame(data []byte,cryptoFrame *CryptoFrame) error {
	rb := util.NewReadBuffer(data)
	for {
		readByte, err := rb.ReadByte()
		if err != nil {
			return err
		}
		if readByte == 0x6 {
			err =cryptoFrame.Read(rb)
			if err != nil {
				return err
			}
		}
		if rb.Buffered() == 0 {
			break
		}
	}
	return nil
}


//func UnCryptoFramePayload(cryptoFrame *CryptoFrame) error {
//	rb := util.NewReadBuffer(cryptoFrame.Data)
//	readByte, err := rb.ReadByte()
//	if err != nil {
//		return err
//	}
//	if tls.HandshakeType(readByte) == tls.ClientHelloType {
//		_, err := tls.ReadClientHello(rb)
//		if err != nil {
//			return err
//		}
//		//log.Print(hello)
//	}
//	return err
//}

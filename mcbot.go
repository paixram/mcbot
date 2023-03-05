package main

import (
	"context"
	"fmt"

	"github.com/liuxsys/mcbot/protocol"
)

func main() {

	varFunc := protocol.NewVar[protocol.VARTYPE](8080, protocol.VARINT)

	varFunc.Process()

	fmt.Printf("%x - %b\n", varFunc.GetDataArray(), varFunc.GetDataArray())
	//fmt.Println(varFunc.GetDataArray())

	longFunc := protocol.NewVar[protocol.VARTYPE](-2147483648, protocol.VARLONG)

	longFunc.Process()

	fmt.Printf("%x - %b\n", longFunc.GetDataArray(), longFunc.GetDataArray())

	pos := protocol.NewPosFromCords(18357644, 831, -20882616)

	fmt.Printf("%b\n", pos.GetData())

	// Protocol LOGS

	fmt.Println("****** PROTOCOL LOGS ******")
	protoconn := protocol.NewConnecionAndBind("mc.universocraft.com", context.Background())

	handshake := protocol.HandShake{
		Proto_version:  protocol.NewVar[protocol.VARTYPE](759, protocol.VARINT),
		Server_address: "mc.minecraft.net",
		Next_state:     protocol.NewVar[protocol.VARTYPE](0x01, protocol.VARINT),
	}

	packet := protocol.Packet{
		PacketID:     0x00,
		PacketLenght: &protocol.VarInt{},
		PacketData:   &handshake,
	}

	protoconn.WritePacket(&packet)
	fmt.Println("Escribiendo el paquete de solicitud de estado...")

	status_request := protocol.Packet{
		PacketID:     0x00,
		PacketLenght: &protocol.VarInt{},
		PacketData:   nil,
	}

	protoconn.WritePacket(&status_request)

	p, _ := protoconn.RecievePacket()

	p.HandlePacket(func(pid uint8, packetLength int64, data []byte) {
		fmt.Printf("\nPacketID: %d - PacketLength: %d - PacketData: \n", pid, packetLength)
	})

	// Iniciar sesion
	handshake_login := protocol.HandShake{
		Proto_version:  protocol.NewVar[protocol.VARTYPE](759, protocol.VARINT),
		Server_address: "mc.minecraft.net",
		Next_state:     protocol.NewVar[protocol.VARTYPE](0x02, protocol.VARINT),
	}

	packet_login := protocol.Packet{
		PacketID:     0x00,
		PacketLenght: &protocol.VarInt{},
		PacketData:   &handshake_login,
	}

	protoconn.WritePacket(&packet_login)

	login_start := protocol.LoginStart{
		Name:          "Lucho19996g",
		HasPlayerUUID: 0x00,
	}

	loginstart_packet := protocol.Packet{
		PacketID:     0x00,
		PacketLenght: &protocol.VarInt{},
		PacketData:   &login_start,
	}

	protoconn.WritePacket(&loginstart_packet)

	protoconn.RecievePacket()

	protoconn.Client.Close()
	//<-finishListenMsg
	fmt.Println("Programa completado")

}

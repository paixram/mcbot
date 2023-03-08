package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/liuxsys/mcbot/protocol"
)

var (
	MULTI_THREADING sync.WaitGroup
	NEXT            chan int = make(chan int, 100)
)

func bot(id_name string) {

	/*fmt.Println("******** VARINT PROTO ********")
	vg := protocol.NewVar[protocol.VARTYPE](1313178, protocol.VARINT)
	vg.Process()
	fmt.Printf("Byte: %x", vg.GetDataArray())

	// Protocol LOGS
	//var SIGNALER chan int
	//SIGNALER = make(chan int, 0x01)
	fmt.Println("****** PROTOCOL LOGS ******")
	protoconn := protocol.NewConnecionAndBind("lacockvid.aternos.me", context.Background()) // mc.deluxezone.net top.mc-complex.com

	// Auto handle Internals (compression set, Reset connections, 0 bytes data received)

	go protoconn.Auto_handler_internals()

	// Iniciar sesion
	fmt.Println("Logs de inicio de session")
	handshake_login := protocol.HandShake{
		Proto_version:  protocol.NewVar[protocol.VARTYPE](759, protocol.VARINT),
		Server_address: "lacockvid.aternos.me",
		Next_state:     protocol.NewVar[protocol.VARTYPE](0x02, protocol.VARINT),
	}

	packet_login := protocol.Packet{
		PacketID:     0x00,
		PacketLenght: &protocol.VarInt{},
		PacketData:   &handshake_login,
	}

	protoconn.WritePacket(packet_login)
	NEXT <- 0x01

	login_start := protocol.LoginStart{
		Name:          id_name,
		HasPlayerUUID: 0x00,
		PlayerUUID:    "38b231576a2a40a7b78cd999dfbb3d50",
	}

	loginstart_packet := protocol.Packet{
		PacketID:     0x00,
		PacketLenght: &protocol.VarInt{},
		PacketData:   &login_start,
	}
	//NEXT <- 0x01
	protoconn.WritePacket(loginstart_packet)
	time.Sleep(1 * time.Second)

	raw_packet, errRawP := protoconn.ReceivePacket() // TODO: Errors in the Packet recieve
	if errRawP != nil {
		fmt.Println("Ocurrio un error", errRawP)

	}
	raw_packet.HandlePacket(func(pid uint8, packetLength int64, data []byte) {
		fmt.Printf("\nPacketID: %d - PacketLength: %d - PacketData:\n", pid, packetLength)
	})

	protoconn.Client.Close()
	close(protoconn.SIGNALER)*/

	//MULTI_THREADING.Done()

	//<-SIGNALER
}

func main() {

	/*for i := 0; i <= 19; i++ {
		name_id := fmt.Sprintf("bot_%d", i)
		MULTI_THREADING.Add(1)
		go bot(name_id)
		<-NEXT
	}

	MULTI_THREADING.Wait()*/

	fmt.Println("******** VARINT PROTO ********")
	vg := protocol.NewVar[protocol.VARTYPE](1313178, protocol.VARINT)
	vg.Process()
	fmt.Printf("Byte: %x", vg.GetDataArray())

	// Protocol LOGS
	//var SIGNALER chan int
	//SIGNALER = make(chan int, 0x01)
	fmt.Println("****** PROTOCOL LOGS ******")
	protoconn := protocol.NewConnecionAndBind("aternos_server", context.Background()) // mc.deluxezone.net top.mc-complex.com

	// Auto handle Internals (compression set, Reset connections, 0 bytes data received)

	go protoconn.Auto_handler_internals()

	// Iniciar sesion
	fmt.Println("Logs de inicio de session")
	handshake_login := protocol.HandShake{
		Proto_version:  protocol.NewVar[protocol.VARTYPE](759, protocol.VARINT),
		Server_address: "aternos_server",
		Next_state:     protocol.NewVar[protocol.VARTYPE](0x02, protocol.VARINT),
	}

	packet_login := protocol.Packet{
		PacketID:     0x00,
		PacketLenght: &protocol.VarInt{},
		PacketData:   &handshake_login,
	}

	protoconn.WritePacket(packet_login)

	login_start := protocol.LoginStart{
		Name:          "dollar_bot",
		HasPlayerUUID: 0x00,
		PlayerUUID:    "38b231576a2a40a7b78cd999dfbb3d50",
	}

	loginstart_packet := protocol.Packet{
		PacketID:     0x00,
		PacketLenght: &protocol.VarInt{},
		PacketData:   &login_start,
	}
	protoconn.WritePacket(loginstart_packet)
	time.Sleep(1 * time.Second)

	raw_packet, errRawP := protoconn.ReceivePacket() // TODO: Errors in the Packet recieve
	if errRawP != nil {
		fmt.Println("Ocurrio un error", errRawP)

	}
	raw_packet.HandlePacket(func(pid uint8, packetLength int64, data []byte) {
		fmt.Printf("\nPacketID: %d - PacketLength: %d - PacketData:\n", pid, packetLength)
	})

	protoconn.Client.Close()
	close(protoconn.SIGNALER)

	fmt.Println("Programa completado")

}

package protocol

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
)

// Tipos del protocol
type long int64
type status int8

// Constantes y mascaras para manipular los bits

const (
	Done  status = 1
	Error status = 0
)

const (
	segment_bits int = 0x7f
	continue_bit int = 0x80
)

// VarInt
type VarInt struct {
	value int
	data  []byte
}

func (vi *VarInt) WriteByte(bv byte) {
	vi.data = append(vi.data, bv)
}

func (vi *VarInt) Process() {
	for {
		if (vi.value & ^segment_bits) == 0 {
			vi.WriteByte(byte(vi.value))
			break
		}

		vi.WriteByte(byte((vi.value & segment_bits) | continue_bit))

		vi.value = int(uint32(vi.value) >> 7)

	}
}

func (vi *VarInt) GetDataArray() []byte {
	return vi.data
}

// VarLong
type VarLong struct {
	value long
	data  []byte
}

func (vl *VarLong) WriteByte(bv byte) {
	vl.data = append(vl.data, bv)
}

func (vl *VarLong) Process() {
	for {
		if (vl.value & ^long(segment_bits)) == 0 {
			vl.WriteByte(byte(vl.value))
			break
		}

		vl.WriteByte(byte((vl.value & long(segment_bits)) | long(continue_bit)))

		vl.value = long(uint64(vl.value) >> 7)
	}
}

func (vl *VarLong) GetDataArray() []byte {
	return vl.data
}

// Declarar y codificar la parte de la creacion de una variable
type VARTYPE int8

const (
	VARINT  VARTYPE = 0x01
	VARLONG VARTYPE = 0x02
)

func NewVar[T VARTYPE](value int64, var_type T) varFuncs {

	switch var_type {
	case T(VARINT):
		fmt.Println("Es de tipo VARTINT")

		return &VarInt{
			value: int(value),
			data:  make([]byte, 0),
		}

	case T(VARLONG):
		fmt.Println("Es de tipo VARLONG")

		return &VarLong{
			value: long(value),
			data:  make([]byte, 0),
		}
	}

	return nil
}

type varFuncs interface {
	WriteByte(bv byte)
	Process()
	GetDataArray() []byte
}

// Definir el tipo de datos de posicion

const (
	bits26Mask = 0x3FFFFFF
	bits12Mask = 0xFFF
)

type Position struct {
	x    int64
	y    int64
	z    int64
	data int64
}

func (p *Position) applyMask() status {
	p.x = (p.x & bits26Mask)
	p.y = (p.y & bits12Mask)
	p.z = (p.z & bits26Mask)

	return Done
}

func (p *Position) process() int64 {
	//fmt.Printf("Position Byte: %b - %b - %b", p.x, p.y, p.z)

	codec := ((p.x & bits26Mask) << 38) | ((p.z & bits26Mask) << 12) | (p.y & bits12Mask)

	fmt.Printf("Total: %b\n", codec)
	return codec
}

func PosReader(data int64) *Position {
	return &Position{
		x: data >> 38,
		y: data << 52 >> 52,
		z: data << 26 >> 38,
	}
}

func (p *Position) GetData() int64 {
	return p.data
}

func NewPosFromCords(x int64, y int64, z int64) *Position {
	position := &Position{
		x: x,
		y: y,
		z: z,
	}

	apply_res := position.applyMask()

	if apply_res != Done {
		log.Fatal("[ + ] No se pudo aplicar las mascaras para las coordenadas")
	}

	data := position.process()

	position.data = data

	return position
}

// Run protocol

/*
type HandShake struct {
	IDpacket       uint8
	proto_version  VarInt
	server_address string
	server_port    uint16
	next_state     VarInt
}*/

// Tipo de datos de paquete

type HandShake struct {
	Proto_version  varFuncs //VarInt
	Server_address string
	Server_port    uint16
	Next_state     varFuncs //VarInt
	Data           []byte
}

func (hs *HandShake) Writer(w io.ByteWriter) {
	for i := 0; i < len(hs.Data); i++ {
		w.WriteByte(hs.Data[i])
	}
}

func (hs *HandShake) packetLen() int {
	//payloadLen := uint32(len(hs.proto_version.data) + len(hs.server_address) + len(hs.next_state.data) + 3)
	hs.Proto_version.Process()
	hs.Next_state.Process()
	fmt.Printf("Proto version: %x - Len of proto_version: %d", hs.Proto_version.GetDataArray(), len(hs.Proto_version.GetDataArray()))
	//fmt.Printf("Nex state: %x - Len of next_state: %d", hs.Next_state.GetDataArray(), len(hs.Next_state.GetDataArray()))

	if len(hs.Data) != 0 {
		return len(hs.Data)
	}

	// Convert PORT to int
	port_buf := make([]byte, 2)
	binary.BigEndian.PutUint16(port_buf, 25565)
	fmt.Printf("Port buf: %x - Port size: %d\n", port_buf, len(port_buf))

	hs.Data = make([]byte, 0)

	hs.Data = append(hs.Data, hs.Proto_version.GetDataArray()...)
	//hs.Data = append(hs.Data, []byte(hs.Server_address)...)
	// Process server addres before append into buffer
	server_encode := EncodeString(hs.Server_address)
	//fmt.Printf("Encode: %x", server_encode)
	hs.Data = append(hs.Data, server_encode...)

	hs.Data = append(hs.Data, port_buf...)
	hs.Data = append(hs.Data, hs.Next_state.GetDataArray()...)

	return len(hs.Data)
}

type PacketIface interface { // All packets data field have this method
	Writer(w io.ByteWriter)
	packetLen() int
}

type Packet struct {
	PacketLenght varFuncs //VarInt
	PacketID     uint8
	PacketData   PacketIface // Todo tipo de paquete contiene el tipo PacketIFace (sus funciones)
}

func (p *Packet) Writer(w io.ByteWriter) {

	var size_packet_without_packetid_field int = 0
	if p.PacketData != nil {
		size_packet_without_packetid_field = p.PacketData.packetLen()
	}

	fmt.Println("Packet Data size: ", size_packet_without_packetid_field)

	//SizePacketLenght := size_packet_without_packetid_field + int(p.PacketID)

	// Procesar el packetID y convertirlo en VarInt
	packetidFace := NewVar[VARTYPE](int64(p.PacketID), VARINT)
	packetidFace.Process()

	packetid_data := packetidFace.GetDataArray()
	fmt.Printf("Packet id data: %x", packetid_data)

	// Procesar el packetLenght y escribirlo en el paquete en formato VarInt

	SizePacketLenght := size_packet_without_packetid_field + len(packetid_data)

	packet_data_format := NewVar[VARTYPE](int64(SizePacketLenght), VARINT)
	packet_data_format.Process()
	packet_data_format_bytes := packet_data_format.GetDataArray()
	fmt.Printf("Packet total lenght: %d - Packet data size content in Varint: %x\n", SizePacketLenght, packet_data_format_bytes)
	for i := 0; i < len(packet_data_format_bytes); i++ {
		w.WriteByte(packet_data_format_bytes[i])
	}

	//w.WriteByte(packetidFace.GetDataArray()...) // Se escribe el ID del paquete
	for i := 0; i < len(packetid_data); i++ {
		w.WriteByte(packetid_data[i]) // Escribir el packetId en formato VarInt en el paquete
	}

	if p.PacketData != nil {
		p.PacketData.Writer(w)
	}

}

// Protocol specs
type ConnectionHandler struct {
	Client net.Conn
	ctx    context.Context
}

func (ch *ConnectionHandler) WritePacket(packet *Packet) status {
	//number_bits_written, errWrite := ch.Client.Write(data)

	//buf := bytes.NewBuffer([]byte{})

	buf := new(bytes.Buffer)
	packet.Writer(buf)

	//fmt.Printf("Buf Size: %d - Buf Data: %b", len(buf.Bytes()), buf.Bytes())

	number_bits_written, errWrite := ch.Client.Write(buf.Bytes()) // buf.Bytes()

	if errWrite != nil || number_bits_written != len(buf.Bytes()) {
		log.Fatal("[ - ] No se pudieron escribir los bytes en en el servidor: ", errWrite)
		return Error
	}

	fmt.Printf("[ + ] Fueron escritos %d Bytes: %x\n", number_bits_written, buf.Bytes())

	return Done

}

/*
	PING START
*/

type Ping struct {
	Payload uint64
	Data    []byte
}

func (p *Ping) Writer(w io.ByteWriter) {
	for i := 0; i < len(p.Data); i++ {
		w.WriteByte(p.Data[i])
	}
}

func (p *Ping) packetLen() int {
	// Procesar el numero payload

	p.Data = append(p.Data, byte(p.Payload))

	return len(p.Data)
}

func HeredatePing(c net.Conn) {
	c.Write([]byte{0xFE, 0x01, 0xFA})
}

// END PING

// ************************************************

/*
	NO PREMIUM LOGIN START
*/

type NPREMIUMLG struct {
	Username string
	Pwd      string
	Data     []byte
}

func (np *NPREMIUMLG) Writer(w io.ByteWriter) {
	for i := 0; i < len(np.Data); i++ {
		w.WriteByte(np.Data[i])
	}
}

func (np *NPREMIUMLG) packetLen() int {
	np.Data = make([]byte, 0)
	user_buf := EncodeString(np.Username)
	pwd_buf := EncodeString(np.Pwd)

	np.Data = append(np.Data, user_buf...)
	np.Data = append(np.Data, pwd_buf...)
	return len(np.Data)
}

// END LOGIN

/*
	LOGIN START
*/

type LoginStart struct {
	Name          string
	HasPlayerUUID uint8
	PlayerUUID    string //big.Int = The real type is bit-128
	Data          []byte
}
type xd big.Int

func (ls *LoginStart) Writer(w io.ByteWriter) {
	for i := 0; i < len(ls.Data); i++ {
		w.WriteByte(ls.Data[i])
	}
}

func (ls *LoginStart) packetLen() int {
	// Procesar el string

	if len(ls.Data) != 0 {
		return len(ls.Data)
	}

	name_databuf := EncodeString(ls.Name)

	ls.Data = append(ls.Data, name_databuf...)
	ls.Data = append(ls.Data, ls.HasPlayerUUID)

	if ls.HasPlayerUUID == 0x01 {
		msb, lsb := EncodeUUID(ls.PlayerUUID)
		msb_buf := make([]byte, 8)
		binary.BigEndian.PutUint64(msb_buf, msb)
		lsb_buf := make([]byte, 8)
		binary.BigEndian.PutUint64(lsb_buf, lsb)
		ls.Data = append(ls.Data, msb_buf...)
		ls.Data = append(ls.Data, lsb_buf...)
		//ls.Data = append(ls.Data, ls.PlayerUUID.Bytes()...)
	}

	return len(ls.Data)
}

func EncodeUUID(UUID string) (uint64, uint64) {
	fmt.Println("************* EMPEZANDO EL ENCODE ******************")
	bytes, _ := hex.DecodeString(UUID)

	// Creando dos enteros sin firmar de 64 bits
	msb := uint64(bytes[0])<<56 | uint64(bytes[1])<<48 | uint64(bytes[2])<<40 | uint64(bytes[3])<<32 | uint64(bytes[4])<<24 | uint64(bytes[5])<<16 | uint64(bytes[6])<<8 | uint64(bytes[7])
	lsb := uint64(bytes[8])<<56 | uint64(bytes[9])<<48 | uint64(bytes[10])<<40 | uint64(bytes[11])<<32 | uint64(bytes[12])<<24 | uint64(bytes[13])<<16 | uint64(bytes[14])<<8 | uint64(bytes[15])

	// Crear el entero sin signo de 128 bits
	//var result big.Int
	//result.Lsh(big.NewInt(0).SetUint64(msb), 64)
	//result.Or(&result, big.NewInt(0).SetUint64(lsb))

	fmt.Printf("Dataa: %x - %x", msb, lsb)
	//return result.Bytes()
	return msb, lsb
}

// END LOGIN

// Manejar los paquetes que llegan
type Handler func(uint8, int64, []byte)
type RecPacketFormat struct {
	PacketLength  int64
	PacketID      uint8
	Data          []byte
	total_payload []byte
}

func (rpf *RecPacketFormat) HandlePacket(h Handler) {
	// Decodificar paquete

	// Get the length of the data and set how many leading bytes of the string to read (2 or 1) MAX_BYTES_LENGTH_FIELD = 2
	var get_n_bytes uint8 = 1
	get_tpayload_len := len(rpf.total_payload)

	if get_tpayload_len > 127 {
		get_n_bytes = 2
	}

	//packet_lenght := rpf.total_payload[get_n_bytes]
	var packet_lenght []byte
	var num int
	if get_tpayload_len > 0 {
		for i := 0; i < int(get_n_bytes); i++ {
			packet_lenght = append(packet_lenght, rpf.total_payload[i])
		}
		num, _ = Undo(packet_lenght)
	}

	rpf.PacketLength = int64(num)
	rpf.PacketID = rpf.total_payload[get_n_bytes]

	if len(rpf.total_payload[get_n_bytes+1:]) > 0 {
		rpf.Data = append(rpf.Data, rpf.total_payload[get_n_bytes+1:]...)
	} else {
		rpf.Data = append(rpf.Data, []byte{0x00, 0x00, 0x00}...)
	}

	// Al ultimo, para ejecutar el Handler con los datos decodificados y proporcionados al usuario
	h(rpf.PacketID, rpf.PacketLength, rpf.Data)
}

func Undo(data []byte) (int, error) {
	value := 0
	position := 0
	counter := 0
	var currentByte byte

	for {
		currentByte = data[counter]
		value |= (int(currentByte) & segment_bits) << position

		if (int(currentByte) & continue_bit) == 0 {
			break
		}

		position += 7
		counter += 1

		if position >= 32 {
			log.Fatal("VarInt es mas grande que 32bits")
		}
	}
	return value, nil
}

type Raw_Packet interface {
	HandlePacket(h Handler)
}

func (ch *ConnectionHandler) ReceivePacket() (Raw_Packet, error) {
	packet_recv := &RecPacketFormat{}
	// Instanciar un reader
	data := new(bytes.Buffer)
	raw_data := bufio.NewWriter(data)

	//n, err := ch.Client.Read()

	n, err := io.Copy(raw_data, ch.Client)

	if err != nil {
		log.Fatal("Se obtuvo un error al leer mensaje: ", err)
	}

	fmt.Printf("[ + ] Llego un mensaje de datos: %x de peso %d Bytes", data.Bytes(), n)

	packet_recv.total_payload = append(packet_recv.total_payload, data.Bytes()...)

	data.Reset()
	return packet_recv, nil

}

func EncodeString(value string) []byte {
	value_length := len(value)

	value_length_varint := NewVar[VARTYPE](int64(value_length), VARINT)

	value_length_varint.Process()

	data_string := append(value_length_varint.GetDataArray(), []byte(value)...)

	return data_string
}

const (
	PORT string = "25565"
)

func NewConnecionAndBind(address string, contexto context.Context) *ConnectionHandler {
	client, errNeti := net.Dial("tcp", address+":"+PORT)

	if errNeti != nil {
		log.Fatal("[ + ] Ha ocurrido un error al conectar y bindear")
		return nil
	}

	fmt.Printf("[ + ] Se ha establecido la conexion con: %s", address)

	return &ConnectionHandler{
		Client: client,
		ctx:    contexto,
	}
}

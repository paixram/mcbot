# MOTIVATION
The motivation for wanting to create a minecraft bot with integrated AI is to learn technical aspects of networking and familiarize myself with network protocols, and what better than to get into building the Minecraft protocol from scratch.
Besides that, I plan to build a bot that can connect to different servers and can be managed by console, for example: "set pvp mode" indicates that the bot will be activated to kill any entity that is not the owner of the bot, either in pvp mode, skywar mode or bedwar mode or even survival. All with the help of AI.

-----
# GOALS AND TODOS
- Achieve a stable connection with servers (NOT PREMIUM)/Cracked for the moment.
  
- Find out about the different forms of authentication implemented by servers such as universocraft.

- Improve efficiency when writing bytes to buffers.

- Fix the Reconnect and resend packet responde log (NOT IS A PROBLEM, IS A LOG)

------


# UPDATES (7/3/2023)
- Reconnection in case of server disconnection.
- Stable connection to the BongeeCord server.
- Handling of packets received with the function: 
```go
p, _ := connection.ReceivePacket()
p.HandlePacket(func(pid uint8, packetLength int64, data []byte) {
	fmt.Printf("PacketID: %d - PacketLength: %d - PacketData: %x", pid, packetLength, data))
    /// Handle Packet...
}
```
---- 

# IMPORTANT!
- ```The default server is one of aternos, if you want to play on another one, you must configure the port in the protocol.go file, set it to 25565.```
  
- ```The bot is only available for connection, I'm still working on the interaction through commands and improving some features that make the protocol a little inefficient.```
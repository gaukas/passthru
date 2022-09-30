## Connection Handler

Connection Handler listens for incoming connections (`net.Conn`), copy the stream (as `net.Conn`s) and feed them into all protocol filters, get the filtering result and forward the connection to the corresponding target.
// Implementation of a KeyValueServer. Students should write their code in this file.

package pa1

import (
	"DS_PA1/rpcs"
	"fmt"
	"net"
	"bufio"
	"strconv"
)

type keyValueServer struct {
	// TODO: implement this!
	totalConnections int
	socketList []net.Conn
	read chan string
	write chan string
	criticalState chan bool

}

// New creates and returns (but does not start) a new KeyValueServer.
func New() KeyValueServer {
	server := new(keyValueServer)
	//server.totalConnections = 0
	//server .socketList = []bufio.ReadWriter{}
	server.read = make(chan string)
	server.criticalState = make(chan bool)
	server.write = make(chan string)
	return server
}

func (kvs *keyValueServer) StartModel1(port int) error {
	// attempt to listen on port 9999
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	// print error and return if couldn't listen on port 999
	if err != nil {
		fmt.Printf("Couldn't listen on port 9999: %s\n", err)
		return nil
	}

	// listen forever
	fmt.Println("Server has started..")
	
	initDB()
	go func() {
		for {
			//fmt.Println("Connected Clients: ",kvs.totalConnections)
			conn, err := ln.Accept()
			kvs.socketList = append(kvs.socketList,conn)
			kvs.totalConnections++
			//fmt.Println("Client",kvs.totalConnections ,"has been connected...")
			
			// handle successful connections concurrently
			if err != nil {
				fmt.Printf("Couldn't accept a client connection: %s\n", err)
			} else {
				go handleConnection(conn,kvs)
			}
		}
	}()
	return nil
}

func (kvs *keyValueServer) Close() {
	// TODO: implement this!
}

func (kvs *keyValueServer) Count() int {
	// TODO: implement this! for first deadline
	return kvs.totalConnections
}

func (kvs *keyValueServer) StartModel2(port int) error {
	// TODO: implement this!
	//
	// Do not forget to call rpcs.Wrap(...) on your kvs struct before
	// passing it to <sv>.Register(...)
	//
	// Wrap ensures that only the desired methods (RecvGet and RecvPut)
	// are available for RPC access. Other KeyValueServer functions
	// such as Close(), StartModel1(), etc. are forbidden for RPCs.
	//
	// Example: <sv>.Register(rpcs.Wrap(kvs))
	return nil
}

func (kvs *keyValueServer) RecvGet(args *rpcs.GetArgs, reply *rpcs.GetReply) error {
	// TODO: implement this!
	return nil
}

func (kvs *keyValueServer) RecvPut(args *rpcs.PutArgs, reply *rpcs.PutReply) error {
	// TODO: implement this!
	return nil
}

// TODO: add additional methods/functions below!

// handleConnection handles client connections
func  handleConnection(conn net.Conn,kvs *keyValueServer) {
	// clean up once the connection closes
	defer Clean(conn)
	

	// obtain a buffered reader / writer on the connection
	rw := ConnectionToRW(conn)
	for {
		// get client message
		msg, err := rw.ReadString('\n')
		if err != nil {
			fmt.Printf("There was an error reading from a client connection: %s\n", err)
			kvs.totalConnections--
			//fmt.Println("Connected Clients: ",kvs.totalConnections)
			return
		}

		// print client message
		//fmt.Printf("Recieved: '%s' of len %d from client: %v\n", msg[:len(msg)-1], len(msg), conn)
		
		command,key,value := parsingData(msg)
		go func() {
			if(command == "put") {
				kvs.write <- "Put Command Recieved"
			} else if (command == "get") {
				kvs.read <- "Get Command Recieved"
			}
		}()
		
		go readWriteDataBase(key,value,kvs)

	
		
		
	
	


		// echo back the same message to the client
		// _, err = rw.WriteString(msg)
		// if err != nil {
		// 	fmt.Printf("There was an error writing to a client connection: %s\n", err)
		// 	return
		// }
		// err = rw.Flush()
		// if err != nil {
		// 	fmt.Printf("There was an error writing to a client connection: %s\n", err)
		// 	return
		// }
	}
}

// Clean closes a connection
func Clean(conn net.Conn) {
	// clean up connection related data structures and goroutines here
	conn.Close()
}

// ConnectionToRW takes a connection and returns a buffered reader / writer on it
func ConnectionToRW(conn net.Conn) *bufio.ReadWriter {
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
}

func check(word string) int  {
	for index,letter := range word {
		if letter == 44 {
			return index
		}
	}
	return -1
}

func readWriteInDataBase(key string,value []byte,kvs *keyValueServer) {
	select {
		case _ = <- kvs.write:
			put(key,value) // write command
			fmt.Println("Value Put -> Key: ",key," Value: ",string(value))
		case _ = <- kvs.read:
			v := get(key) // read command
			fmt.Println("Value Get -> Key: ",key," Value: ",string(v))
				for _,socket := range kvs.socketList {
				socket.Write([]byte(string(key)+","+string(v)))
			}
	}
}

func readWriteToClient() {
	// get client message
		msg, err := rw.ReadString('\n')
		if err != nil {
			fmt.Printf("There was an error reading from a client connection: %s\n", err)
			kvs.totalConnections--
			//fmt.Println("Connected Clients: ",kvs.totalConnections)
			return
		}
}
func parsingData(msg string) (string,string,[]byte) {
	
	var command,key,value string

	for i:= 0; i <3 ; i ++ {
		spaceIndex := check(msg)
		if command == "get" {
			key = msg[:len(msg)-1]
			break
		} else if i == 0 {
			command = msg[:spaceIndex]
		} else if i == 1 {
			key = msg[:spaceIndex]
		} else {
			value = msg
		}
		msg = msg[spaceIndex+1:]
	}
	return command,key,[]byte(value)
}

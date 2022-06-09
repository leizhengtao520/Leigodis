package parser

import (
	"bufio"
	"errors"
	"io"
)

type Reply interface{
	ToBytes() []byte
}
type Payload struct {
	Data Reply
	Err error
}

func ParseStream(reader io.Reader) <-chan *Payload{
	ch :=make(chan *Payload)
	g
}

func parse0(reader io.Reader,ch chan<-*Payload){







}

type readState struct {
	readingMultiLine bool
	expectedArgsCount int
	msgType byte
	args [][]byte
	bulkLen int64
}

func readLine(bufReader *bufio.Reader,state *readState)([]byte,bool,error ){
	var msg []byte
	var err error
	if state.bulkLen == 0{ //单行
		msg,err = bufReader.ReadBytes('\n')
		if err !=nil{
			return nil, true, err
		}
		if len(msg)==0 || msg[len(msg)-2]!='\r'{
			return nil,false,errors.New("protocol error:"+string(msg))
		}
	}else{//多行
		msg=make([]byte,state.bulkLen+2)// $3
		_,err = io.ReadFull(bufReader,msg)
		if err !=nil{
			return nil,true,err
		}
		if len(msg) ==0 || msg[len(msg)-2] !='\r' || msg[len(msg)-1] != '\n'{
			return nil ,false,errors.New("protocol error:"+string(msg))
		}
		state.bulkLen=0
	}
	return msg, false, err
}




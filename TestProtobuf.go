package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"test"
)

func main() {

	t := &test.MyMsg{}
	t.Str = proto.String("Hi")
	t.Opt = proto.Int32(32)
	t.Id = proto.Int32(32)

	proto.String("Hi")
	proto.Int32(32)

	mdata, err := proto.Marshal(t)

	if err != nil {
		log.Fatal("marshaling error:", err)
	}

	// newTest := &test.MyMsg{}
	var umData test.MyMsg

	err = proto.Unmarshal(mdata, &umData)
	if err != nil {
		log.Fatal("unmarshaling error:", err)
	}

	fmt.Println(umData.GetId())

}

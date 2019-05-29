/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/solution9th/NSBridge/dns_pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	//address = "dns.dev.x3.work:9191"
	address     = "localhost:9191"
	defaultName = "world"
)

func main() {

	creds, err := credentials.NewClientTLSFromFile("../keys/grpc/grpc.pem", "dev")

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	//conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDNSServerClient(conn)

	// Contact the server and print out its response.
	//name := defaultName
	//if len(os.Args) > 1 {
	//	name = os.Args[1]
	//}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	//if err != nil {
	//	log.Fatalf("could not greet: %v", err)
	//}
	//log.Printf("Greeting: %s", r.Message)

	//testDomainList(c, ctx)
	//testDomainCreate(c, ctx)
	//testDomainDelete(c, ctx)
	//testDomainTakeOver(c, ctx)
	// testDomainTypeList(c, ctx)

	//testDomainOfRK(c, ctx)
	//testOwnDomainList(c, ctx)

	//testRecordCreate(c, ctx)
	//testRecordList(c, ctx)
	testRecordDisable(c, ctx)
	//testRecordDelete(c, ctx)
	//testRecordInfo(c, ctx)
	//testRecordUpdate(c, ctx)

}

// test domain
func testDomainList(c pb.DNSServerClient, ctx context.Context) {
	r2, err := c.DomainsList(ctx, &pb.RequestDomainsList{
		Start: 10,
		Count: 99999,
	})
	if err != nil {
		fmt.Println("🌲", err.Error())
	} else {
		fmt.Println("success data:", r2.ErrCode, r2.ErrMsg)
		fmt.Println(GenJson(r2))
	}
}

func testDomainCreate(c pb.DNSServerClient, ctx context.Context) {
	r2, err := c.DomainCreate(ctx, &pb.RequestDomainCreate{
		ApiKey: "f3ed8235-39a2-4577-a2c1-a8e955a38621",
		Domain: fmt.Sprintf("%vmgfeng.com", rand.Intn(2343234)),
		Remark: "测试",
	})
	if err != nil {
		fmt.Println("🌲", err.Error())
	} else {
		fmt.Println(GenJson(r2))
	}
}

func testDomainTakeOver(c pb.DNSServerClient, ctx context.Context) {
	r2, err := c.DomainStatus(ctx, &pb.RequestDomainStatus{
		ApiKey: "cmVjb3Jk2cab9d2edfd19538",
	})
	if err != nil {
		fmt.Println("🌲", err.Error())
	} else {
		fmt.Println(GenJson(r2))
	}
}

func testDomainOfRK(c pb.DNSServerClient, ctx context.Context) {
	r2, err := c.RecordDomainOfRK(ctx, &pb.RequestDomainOfRK{
		RecordKey: "cmVjb3Jk2cab9d2edfd1958",
	})
	if err != nil {
		fmt.Println("🌲", err.Error())
	} else {
		fmt.Println(GenJson(r2))
	}
}

func testOwnDomainList(c pb.DNSServerClient, ctx context.Context) {
	r2, err := c.OwnDomainsList(ctx, &pb.RequestOwnDomainsList{
		//Key: "e00de6c8-79fa-4721-94cf-7311e5f1c9fb",
		Key:"1e1b73bd-f817-4dbe-9991-89a2ec05d7f4",
	})

	if err != nil {
		fmt.Println("🌲", err.Error())
	} else {
		fmt.Println(GenJson(r2))
	}
}

func testDomainDelete(c pb.DNSServerClient, ctx context.Context) {
	r2, err := c.DomainDelete(ctx, &pb.RequestDomainDelete{
		ApiKey: "e00de6c8-79fa-4721-94cf-7311e5f1c9fb",
		Id:     66,
	})
	if err != nil {
		fmt.Println("🌲", err.Error())
	} else {
		fmt.Println(GenJson(r2))
	}
}

func testDomainTypeList(c pb.DNSServerClient, ctx context.Context) {
	r2, err := c.Types(ctx, &empty.Empty{})
	if err != nil {
		fmt.Println("🌲", err.Error())
	} else {
		fmt.Println(GenJson(r2))
	}
}

// test record
func testRecordList(c pb.DNSServerClient, ctx context.Context) {
	// test record operation
	r1, err := c.RecordList(ctx, &pb.RequestRecordList{
		RecordKey: "cmVjb3Jk2cab9d2edfd19539",
	})
	if err != nil {
		fmt.Println("🌧", err.Error())
	} else {
		fmt.Println(GenJson(r1))
	}
}

func testRecordInfo(c pb.DNSServerClient, ctx context.Context) {
	// test record operation
	r1, err := c.RecordInfo(ctx, &pb.RequestRecordInfo{
		RecordKey: "cmVjb3Jk2cab9d2edfd19539",
		RecordId:  18,
	})
	if err != nil {
		fmt.Println("🌧", err.Error())
	} else {
		fmt.Println(GenJson(r1))
	}
}

func testRecordUpdate(c pb.DNSServerClient, ctx context.Context) {
	// test record operation
	r1, err := c.RecordUpdate(ctx, &pb.RequestRecordUpdate{
		RecordKey:  "cmVjb3Jk2cab9d2edfd19539",
		SubDomain:  "uus2u.com",
		RecordType: "NS",
		Unit:       "day",
		Value:      "mafeng.com",
		LineId:     2,
		RecordId:   37,
	})
	if err != nil {
		fmt.Println("🌧", err.Error())
	} else {
		fmt.Println(GenJson(r1))
	}
}

func testRecordDisable(c pb.DNSServerClient, ctx context.Context) {
	// test record operation
	r1, err := c.RecordDisable(ctx, &pb.RequestRecordDisable{
		RecordKey: "cmVjb3Jk2cab9d2edfd19539",
		RecordId:  18,
		Disable:   false,
	})
	if err != nil {
		fmt.Println("🌧", err.Error())
	} else {
		fmt.Println(GenJson(r1), r1.ErrCode, r1.ErrMsg, r1)
		fmt.Printf("%#v\n",r1)
	}
}

func testRecordDelete(c pb.DNSServerClient, ctx context.Context) {
	// test record operation
	r1, err := c.RecordDelete(ctx, &pb.RequestRecordDelete{
		RecordKey: "cmVjb3Jk2cab9d2edfd19539",
		RecordId:  13,
	})
	if err != nil {
		fmt.Println("🌧", err.Error())
	} else {
		fmt.Println(GenJson(r1))
	}
}

func testRecordCreate(c pb.DNSServerClient, ctx context.Context) {
	// test record operation
	rand.Seed(int64(time.Now().Nanosecond()))
	r1, err := c.RecordCreate(ctx, &pb.RequestRecordCreate{
		RecordKey:  "cmVjb3Jk09abadfee73490e9",
		SubDomain:  fmt.Sprintf("%vmgfeng.com", rand.Intn(43234)),
		RecordType: "TXT",
		Unit:       "day",
		Value:      fmt.Sprintf("%vvalue.com", rand.Intn(43234)),
		LineId:     int64(rand.Intn(500)),
	})
	if err != nil {
		fmt.Println("🌧", err.Error())
	} else {
		fmt.Println(GenJson(r1))
	}
}

func GenJson(v interface{}) string {
	js, _ := json.Marshal(v)
	return string(js)
}

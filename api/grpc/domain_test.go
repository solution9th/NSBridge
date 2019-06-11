package grpc

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/solution9th/NSBridge/dns_pb"
	"github.com/solution9th/NSBridge/internal/utils"
)

var (
	ctx = context.Background()
)

func TestCreateDomain(t *testing.T) {

	InitDB(t)
	c := New()

	r, err := c.DomainCreate(ctx, &pb.RequestDomainCreate{
		ApiKey: "123",
		Domain: "u22us2u.com",
		Remark: "from grpc",
	})
	if err != nil {
		t.Error("create domain error:", err)
	}

	fmt.Println(utils.GenJson(r))

}

func TestListDomains(t *testing.T) {
	InitDB(t)

	c := New()

	tests := []struct {
		req  *pb.RequestDomainsList
		code int
		err  error
	}{
		{&pb.RequestDomainsList{}, 0, nil},
		{&pb.RequestDomainsList{
			Start: 10,
			Count: 99999,
		}, 0, nil},
	}

	for k, v := range tests {

		resp, err := c.DomainsList(ctx, v.req)
		if err != v.err {
			t.Errorf("err k:%v want: %v,got: %v", k, v.err, err)
		}

		if int(resp.ErrCode) != v.code {
			t.Errorf("code k: %v, code: %v,want: %v", k, resp.ErrCode, v.code)
		}
	}

}

func TestStatusDomain(t *testing.T) {

	InitDB(t)

	c := New()

	r, err := c.DomainStatus(ctx, &pb.RequestDomainStatus{
		ApiKey: "cmVjb3Jk2cab9d2edfd19538",
	})

	if err != nil {
		t.Error("domain status error:", err)
	}

	fmt.Println(err, utils.GenJson(r))

}

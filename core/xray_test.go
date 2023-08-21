package core_test

import (
	"testing"

	"github.com/AikoPanel/Xray-core/app/dispatcher"
	"github.com/AikoPanel/Xray-core/app/proxyman"
	"github.com/AikoPanel/Xray-core/common"
	"github.com/AikoPanel/Xray-core/common/net"
	"github.com/AikoPanel/Xray-core/common/protocol"
	"github.com/AikoPanel/Xray-core/common/serial"
	"github.com/AikoPanel/Xray-core/common/uuid"
	. "github.com/AikoPanel/Xray-core/core"
	"github.com/AikoPanel/Xray-core/features/dns"
	"github.com/AikoPanel/Xray-core/features/dns/localdns"
	_ "github.com/AikoPanel/Xray-core/main/distro/all"
	"github.com/AikoPanel/Xray-core/proxy/dokodemo"
	"github.com/AikoPanel/Xray-core/proxy/vmess"
	"github.com/AikoPanel/Xray-core/proxy/vmess/outbound"
	"github.com/AikoPanel/Xray-core/testing/servers/tcp"
	"google.golang.org/protobuf/proto"
)

func TestXrayDependency(t *testing.T) {
	instance := new(Instance)

	wait := make(chan bool, 1)
	instance.RequireFeatures(func(d dns.Client) {
		if d == nil {
			t.Error("expected dns client fulfilled, but actually nil")
		}
		wait <- true
	})
	instance.AddFeature(localdns.New())
	<-wait
}

func TestXrayClose(t *testing.T) {
	port := tcp.PickPort()

	userID := uuid.New()
	config := &Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
		Inbound: []*InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortList: &net.PortList{
						Range: []*net.PortRange{net.SinglePortRange(port)},
					},
					Listen: net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
					Address: net.NewIPOrDomain(net.LocalHostIP),
					Port:    uint32(0),
					NetworkList: &net.NetworkList{
						Network: []net.Network{net.Network_TCP},
					},
				}),
			},
		},
		Outbound: []*OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&outbound.Config{
					Receiver: []*protocol.ServerEndpoint{
						{
							Address: net.NewIPOrDomain(net.LocalHostIP),
							Port:    uint32(0),
							User: []*protocol.User{
								{
									Account: serial.ToTypedMessage(&vmess.Account{
										Id: userID.String(),
									}),
								},
							},
						},
					},
				}),
			},
		},
	}

	cfgBytes, err := proto.Marshal(config)
	common.Must(err)

	server, err := StartInstance("protobuf", cfgBytes)
	common.Must(err)
	server.Close()
}

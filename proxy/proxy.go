// Package proxy is a cli proxy
package proxy

import (
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config/options"
	"github.com/micro/go-micro/proxy"
	"github.com/micro/go-micro/proxy/mucp"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/util/log"
)

var (
	// Name of the proxy
	Name = "go.micro.proxy"
	// The address of the proxy
	Address = ":8081"
	// The endpoint host to route to
	Endpoint string
)

func run(ctx *cli.Context, srvOpts ...micro.Option) {
	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("endpoint")) > 0 {
		Endpoint = ctx.String("endpoint")
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// service opts
	srvOpts = append(srvOpts, micro.Name(Name))
	if i := time.Duration(ctx.GlobalInt("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(ctx.GlobalInt("register_interval")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterInterval(i*time.Second))
	}

	// set address
	if len(Address) > 0 {
		srvOpts = append(srvOpts, micro.Address(Address))
	}

	// set the context
	var popts []options.Option

	// set endpoint
	if len(Endpoint) > 0 {
		popts = append(popts, proxy.WithEndpoint(Endpoint))
	}

	// new proxy
	p := mucp.NewProxy(popts...)

	// new service
	service := micro.NewService(srvOpts...)

	// set the router
	service.Server().Init(
		server.WithRouter(p),
	)

	if len(Endpoint) > 0 {
		log.Logf("Proxy [%s] Serving endpoint %s\n", p.String(), Endpoint)
	}

	// Run internal service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "proxy",
		Usage: "Run the service proxy",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the proxy http address e.g 0.0.0.0:8081",
				EnvVar: "MICRO_PROXY_ADDRESS",
			},
			cli.StringFlag{
				Name:   "endpoint",
				Usage:  "Set the endpoint to route to e.g greeter or localhost:9090",
				EnvVar: "MICRO_PROXY_ENDPOINT",
			},
		},
		Action: func(ctx *cli.Context) {
			run(ctx, options...)
		},
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	return []cli.Command{command}
}

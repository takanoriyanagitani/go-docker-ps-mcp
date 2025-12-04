package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/takanoriyanagitani/go-docker-ps-mcp/ps"
)

type PsInput struct{}

type PsOutput struct {
	Result []ps.BasicSummary `json:"result"`
}

var dockerHost = flag.String("docker-host", "/var/run/docker.sock", "Path to the Docker Unix socket")
var port = flag.Int("port", 12028, "port to listen")

var containerNamePattern = flag.String("container-name-pattern", "^/cadv.*$", "Container name pattern")

func PsTool(ctx context.Context, req *mcp.CallToolRequest, input PsInput) (
	*mcp.CallToolResult,
	PsOutput,
	error,
) {
	containers, err := ps.GetContainers(ctx, *dockerHost)
	if err != nil {
		return nil, PsOutput{}, err
	}

	var patString ps.ContainerNamePatternString = ps.ContainerNamePatternString(*containerNamePattern)
	filter, e := patString.ToSummaryFilter()
	if nil != e {
		return nil, PsOutput{}, e
	}

	var basicList []ps.BasicSummary = ps.
		SummaryList(containers).
		ToBasicList(filter)

	return nil, PsOutput{Result: basicList}, nil
}

func main() {
	flag.Parse()
	var server *mcp.Server = mcp.NewServer(&mcp.Implementation{Name: "docker-ps", Version: "v1.0.0"}, nil)
	var mhandler *mcp.StreamableHTTPHandler = mcp.NewStreamableHTTPHandler(
		func(req *http.Request) *mcp.Server { return server },
		&mcp.StreamableHTTPOptions{Stateless: true},
	)
	var handler http.Handler = mhandler
	var address string = fmt.Sprintf(":%d", *port)
	hserver := &http.Server{
		Addr:           address,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	mcp.AddTool(server, &mcp.Tool{Name: "ps", Description: "list docker containers"}, PsTool)
	log.Printf("ready to start http mcp server. use %s\n", address)
	hse := hserver.ListenAndServe()
	if nil != hse {
		log.Fatalf("%v\n", hse)
	}
}

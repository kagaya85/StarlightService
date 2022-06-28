package lb

import (
	"context"
	"os"
	"sync"

	v1 "../starlight/api/balancer/v1"

	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type Weight struct {
	instance string
	weight   int
}

type WeightList struct {
	mu   sync.Mutex
	list map[string][]Weight // key: operation, value: weight list
}

func (w *WeightList) Sync(ctx context.Context, conf *conf.Balancer) error {
	for {
		conn, err := grpc.Dial(ctx, grpc.WithEndpoint(conf.Addr))
		if err != nil {
			return err
		}
		client := v1.NewWeightUpdaterClient(conn)

	}
	return nil
}

func listInstanceInfo() (ins, pod, node, zone string) {
	ins = os.Getenv("POD_IP")
	pod = os.Getenv("POD_NAME")
	node = os.Getenv("NODE_NAME")
	zone = os.Getenv("ZONE_NAME")
	return
}

func listServiceInfo() (services []*v1.ServiceInfo, err error) {
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Services().Len(); i++ {
			var (
				methods []string
				sd      = fd.Services().Get(i)
			)
			for j := 0; j < sd.Methods().Len(); j++ {
				md := sd.Methods().Get(j)
				mName := string(md.Name())
				methods = append(methods, mName)
			}
			services = append(services, &v1.ServiceInfo{
				Service:    string(sd.Name()),
				port:       0,
				Operations: methods,
			})
		}
		return true
	})
	return
}

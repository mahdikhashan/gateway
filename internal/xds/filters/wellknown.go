// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package filters

import (
	grpcstats "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/grpc_stats/v3"
	grpcweb "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/grpc_web/v3"
	httprouter "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/envoyproxy/gateway/internal/utils/protocov"
)

var GRPCWeb, GRPCStats *hcm.HttpFilter

func init() {
	any, err := protocov.ToAnyWithValidation(&grpcweb.GrpcWeb{})
	if err != nil {
		panic(err)
	}
	GRPCWeb = &hcm.HttpFilter{
		Name: wellknown.GRPCWeb,
		ConfigType: &hcm.HttpFilter_TypedConfig{
			TypedConfig: any,
		},
	}

	any, err = protocov.ToAnyWithValidation(&grpcstats.FilterConfig{
		EmitFilterState: true,
		PerMethodStatSpecifier: &grpcstats.FilterConfig_StatsForAllMethods{
			StatsForAllMethods: &wrapperspb.BoolValue{Value: true},
		},
	})
	if err != nil {
		panic(err)
	}
	GRPCStats = &hcm.HttpFilter{
		Name: wellknown.HTTPGRPCStats,
		ConfigType: &hcm.HttpFilter_TypedConfig{
			TypedConfig: any,
		},
	}
}

func GenerateRouterFilter(enableEnvoyHeaders bool) (*hcm.HttpFilter, error) {
	any, err := protocov.ToAnyWithValidation(&httprouter.Router{
		SuppressEnvoyHeaders: !enableEnvoyHeaders,
	})
	if err != nil {
		return nil, err
	}
	return &hcm.HttpFilter{
		Name: wellknown.Router,
		ConfigType: &hcm.HttpFilter_TypedConfig{
			TypedConfig: any,
		},
	}, nil
}

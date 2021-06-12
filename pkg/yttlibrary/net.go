// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package yttlibrary

import (
	"errors"
	"fmt"
	"net"

	"github.com/k14s/starlark-go/starlark"
	"github.com/k14s/starlark-go/starlarkstruct"
	"github.com/k14s/ytt/pkg/orderedmap"
	"github.com/k14s/ytt/pkg/template/core"
)

var (
	NETAPI = starlark.StringDict{
		"net": &starlarkstruct.Module{
			Name: "net",
			Members: starlark.StringDict{
				"parse_ip": starlark.NewBuiltin("net.parse_ip", core.ErrWrapper(netModule{}.ParseURL)),
			},
		},
	}
)

type netModule struct{}

// IPValue stores a parsed IP
type IPValue struct {
	ip                   net.IP
	*core.StarlarkStruct // TODO: keep authorship of the interface by delegating instead of embedding
}

var invalidIPErr = errors.New("input IP string is invalid")

func (b netModule) ParseURL(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if args.Len() != 1 {
		return starlark.None, fmt.Errorf("expected exactly one argument")
	}

	ipStr, err := core.NewStarlarkValue(args.Index(0)).AsString()
	if err != nil {
		return starlark.None, err
	}

	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return starlark.None, invalidIPErr
	}

	return (&IPValue{parsedIP, nil}).AsStarlarkValue(), nil
}

func (iv *IPValue) Type() string { return "@ytt:net.ip" }

func (iv *IPValue) AsStarlarkValue() starlark.Value {
	m := orderedmap.NewMap()
	m.Set("is_ipv4", starlark.NewBuiltin("ip.is_ipv4", core.ErrWrapper(iv.IsIPv4)))
	m.Set("is_ipv6", starlark.NewBuiltin("ip.is_ipv6", core.ErrWrapper(iv.IsIPv6)))
	iv.StarlarkStruct = core.NewStarlarkStruct(m)
	return iv
}

// func (uv *IPValue) ConversionHint() string {
// 	return "IPValue does not automatically encode (hint: use .string())"
// }

func (iv *IPValue) IsIPv4(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if args.Len() != 0 {
		return starlark.None, fmt.Errorf("expected no argument")
	}
	isV4 := iv.ip != nil && iv.ip.To4() != nil
	return starlark.Bool(isV4), nil
}

func (iv *IPValue) IsIPv6(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if args.Len() != 0 {
		return starlark.None, fmt.Errorf("expected no argument")
	}
	isV6 := iv.ip != nil && iv.ip.To16() != nil
	return starlark.Bool(isV6), nil
}

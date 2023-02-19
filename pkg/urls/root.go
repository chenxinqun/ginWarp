package urls

import (
	"fmt"
	"strings"
)

type Url interface {
	Root() string
	Service() string
	Version() string
	Group() string
	Path() string
	String() string
}

type UrlDefine struct {
	root    string
	service string
	// api版本, V1, V2...
	version string
	// 路由组, 对应一个handler
	group string
	// 对应handler下的一个方法
	path string
}

func New() *UrlDefine {
	return &UrlDefine{}
}

// R 设置根路径
func (d *UrlDefine) R(r string) *UrlDefine {
	d.root = r
	return d
}

// S 设置 ServiceName
func (d *UrlDefine) S(s string) *UrlDefine {
	d.service = s
	return d
}

// V 设置version
func (d *UrlDefine) V(v string) *UrlDefine {
	d.version = v
	return d //nolint:nolintlint,govet
}

// G 设置group
func (d *UrlDefine) G(g string) *UrlDefine {
	d.group = g
	return d
}

// P 设置path
func (d *UrlDefine) P(p string) *UrlDefine {
	d.path = p
	return d
}

// Root 获取根路径
func (d *UrlDefine) Root() string {
	return d.root
}

// Service 获取ServiceName
func (d *UrlDefine) Service() string {
	return d.service
}

// Version 获取version
func (d *UrlDefine) Version() string {
	return d.version //nolint:nolintlint,govet
}

// Group 获取group
func (d *UrlDefine) Group() string {
	return d.group
}

// Path 获取path
func (d *UrlDefine) Path() string {
	return d.path
}

// String 获取完整的URL
func (d *UrlDefine) String() string {
	Root := strings.TrimSpace(d.root)
	Service := strings.TrimSpace(d.service)
	Version := strings.TrimSpace(d.version)
	Group := strings.TrimSpace(d.group)
	Path := strings.TrimSpace(d.path)
	Root = strings.Trim(Root, "/")
	Service = strings.Trim(Service, "/")
	Version = strings.Trim(Version, "/")
	Group = strings.Trim(Group, "/")
	Path = strings.Trim(Path, "/")
	url := fmt.Sprintf("/%s/%s/%s/%s/%s", Root, Service, Version, Group, Path)
	return url
}

// RootGroup 根路径
const (
	RootGroup = "/api"
)

// 版本相关常量
const (
	V1 = "/v1"
	V2 = "/v2"
	V3 = "/v3"
	V4 = "/v4"
	V5 = "/v5"
	V6 = "/v6"
	V7 = "/v7"
	V8 = "/v8"
	V9 = "/v9"
)

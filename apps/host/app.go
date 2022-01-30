package host

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cunmao-Jazz/cmdb/apps/resource"
	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
)

const (
	AppName = "host"
)

var (
	validate = validator.New()
)

func (h *Host) GenHash() error {
	hash := sha1.New()

	b, err := json.Marshal(h.Information)
	if err != nil {
		return err
	}
	hash.Write(b)
	h.Base.ResourceHash = fmt.Sprintf("%x", hash.Sum(nil))

	b, err = json.Marshal(h.Describe)
	if err != nil {
		return err
	}
	hash.Reset()
	hash.Write(b)
	h.Base.DescribeHash = fmt.Sprintf("%x", hash.Sum(nil))
	return nil
}

func (d *Describe) LoadKeyPairNameString(str string) {
	if str != "" {
		d.KeyPairName = strings.Split(str, ",")
	}
}

func (d *Describe) LoadSecurityGroupsString(str string) {
	if str != "" {
		d.SecurityGroups = strings.Split(str, ",")
	}
}

//主机登录凭证的引用
func (d *Describe) KeyPairNameToString() string {
	return strings.Join(d.KeyPairName, ",")
}

//主机安全组对象
func (d *Describe) SecurityGroupsToString() string {
	return strings.Join(d.SecurityGroups, ",")
}

func (req *DescribeHostRequest) Where() (string, interface{}) {
	switch req.DescribeBy {
	case DescribeBy_HOST_ID:
		return "id = ?", req.Value
	default:
		return "", nil
	}
}

func NewDefaultHost() *Host {
	return &Host{
		Base:        &resource.Base{},
		Information: &resource.Information{Tags: map[string]string{}},
		ReleasePlan: &resource.ReleasePlan{},
		Describe:    &Describe{},
	}
}

func NewDescribeHostRequestById(id string) *DescribeHostRequest {
	return &DescribeHostRequest{
		DescribeBy: DescribeBy_HOST_ID,
		Value:      id,
	}
}

func NewHostSet() *HostSet {
	return &HostSet{
		Items: []*Host{},
	}
}

func (s *HostSet) Add(item *Host) {
	s.Items = append(s.Items, item)
}

//更新请求的字段值校验
func (req *UpdateHostRequest) Validate() error {
	return validate.Struct(req)
}

func (h *Host) Put(req *UpdateHostData) {
	// 提前保留老的hash
	oldRH, oldDH := h.Base.ResourceHash, h.Base.DescribeHash

	// 对象更新完成后 重新计算Hash
	h.Information = req.Information
	h.Describe = req.Describe
	h.Information.UpdateAt = time.Now().UnixMilli()
	h.GenHash()

	// 更新后的Hash和之前的Hash是否一致，来判断该对象是否需要更新
	if h.Base.ResourceHash != oldRH {
		h.Base.ResourceHashChanged = true
	}
	if h.Base.DescribeHash != oldDH {
		h.Base.DescribeHashChanged = true
	}
}

func (h *Host) Patch(req *UpdateHostData) error {
	oldRH, oldDH := h.Base.ResourceHash, h.Base.DescribeHash

	// patch information
	err := mergo.MergeWithOverwrite(h.Information, req.Information)
	if err != nil {
		return err
	}

	err = mergo.MergeWithOverwrite(h.Describe, req.Describe)
	if err != nil {
		return err
	}

	h.Information.UpdateAt = time.Now().UnixMilli()
	h.GenHash()

	if h.Base.ResourceHash != oldRH {
		h.Base.ResourceHashChanged = true
	}
	if h.Base.DescribeHash != oldDH {
		h.Base.DescribeHashChanged = true
	}

	return nil
}

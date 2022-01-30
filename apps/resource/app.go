package resource

import "strings"

const (
	AppName = "resource"
)

//定义resource集合的构造函数
//注意!!!  指针类型和引用类型一定要初始化,不然会造成空指针错误
func NewResourceSet() *ResourceSet {
	return &ResourceSet{
		Items: []*Resource{},
	}
}

//定义resource的构造函数
//注意!!!  指针类型和引用类型一定要初始化,不然会造成空指针错误
func NewDefaultResource() *Resource {
	return &Resource{
		Base:        &Base{},
		Information: &Information{Tags: map[string]string{}},
		ReleasePlan: &ReleasePlan{},
	}
}
//Information方法 将获取出来的ip转换成字符串切片
func (i *Information) LoadIPString(PrivateIp,PublicIp string) {
	i.PrivateIp = strings.Split(PrivateIp,",")
	i.PublicIp = strings.Split(PublicIp,",")

}

func (i *Information) PublicIPToString () (ips string) {
	return strings.Join(i.PublicIp,",")
}

func (i *Information) PrivateIPToString () (ips string) {
	return strings.Join(i.PrivateIp,",")
}

//resource方法,将查询出来的resource 添加到集合中
func (s *ResourceSet) Add (item *Resource) {
	s.Items = append(s.Items, item)
}

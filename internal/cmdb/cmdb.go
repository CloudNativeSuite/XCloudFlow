package cmdb

// Resource 表示需要登记到 CMDB 的基础结构。
type Resource struct {
	ID   string
	Type string
	Tags map[string]string
}

// Register 在资源创建后向 CMDB 注册信息。
func Register(r Resource) error {
	// TODO: 调用实际的 CMDB 接口
	return nil
}

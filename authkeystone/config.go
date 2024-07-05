package authkeystone

// 认证所需的struct
type AuthSetUp struct {
	// keystone 认证端点
	AuthEndpoint string
	// 认证用户
	UserName string
	// 认证用户的密码
	Password string
	// 认证用户所属的域
	DomainName string
	// 认证用户所属的项目名
	TenantName string
}

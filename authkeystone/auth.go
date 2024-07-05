package authkeystone

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

// 实现openstack 的keystone 认证
func (a *AuthSetUp) Auth() (*gophercloud.ProviderClient, error) {
	// 设置认证所需的值
	authOps := gophercloud.AuthOptions{
		IdentityEndpoint: a.AuthEndpoint,
		Username:         a.UserName,
		Password:         a.Password,
		DomainName:       a.DomainName,
		TenantName:       a.TenantName,
	}
	return openstack.AuthenticatedClient(authOps)
}

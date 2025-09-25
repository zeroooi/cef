// Package security 提供应用程序的安全控制功能
// 包含域名白名单验证、URL过滤等安全机制
package security

import (
	"cef/internal/config"
	"fmt"
	"net/url"
	"strings"
)

const (
	AdLoginUrl = "https://ad.oceanengine.com/pages/login/index.html"
)

// WhitelistValidator 白名单验证器
type WhitelistValidator struct {
	config func(...string) *config.WhitelistConfig
}

// NewWhitelistValidator 创建新的白名单验证器实例
func NewWhitelistValidator(cfg func(...string) *config.WhitelistConfig) *WhitelistValidator {
	return &WhitelistValidator{
		config: cfg,
	}
}

// IsURLAllowed 检查URL是否被允许访问
// 支持精确匹配和子域名匹配两种模式
func (v *WhitelistValidator) IsURLAllowed(requestURL string, account ...string) bool {
	// 解析URL
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		fmt.Printf("URL解析失败: %v\n", err)
		return false
	}
	if parsedURL.Scheme == "bytedance" {
		return true
	}
	// 特殊判断AD、千川系统的登录页面,不允许跳转
	if requestURL == AdLoginUrl {
		return false
	}

	// 获取域名并转换为小写
	hostname := strings.ToLower(parsedURL.Hostname())

	// 检查是否在白名单中
	inAllowedDomain := false
	for _, allowedDomain := range v.config(account...).AllowedDomains {
		allowedDomain = strings.ToLower(allowedDomain)

		// 支持精确匹配和子域名匹配
		if hostname == allowedDomain || strings.HasSuffix(hostname, "."+allowedDomain) {
			inAllowedDomain = true
		}
	}
	if !inAllowedDomain {
		return false
	}
	// 检查是否在黑名单中
	for _, notAllowedDomain := range v.config(account...).NotAllowedDomains {
		notAllowedDomain = strings.ToLower(notAllowedDomain)

		// 支持精确匹配和子域名匹配
		if hostname == notAllowedDomain || strings.HasSuffix(hostname, "."+notAllowedDomain) {
			return false
		}
	}

	return true
}

// GetBlockedMessage 获取访问被阻止时的消息
func (v *WhitelistValidator) GetBlockedMessage() string {
	return v.config().BlockedMessage
}

// GetRedirectURL 获取被阻止时的重定向URL
func (v *WhitelistValidator) GetRedirectURL() string {
	return v.config().RedirectBlockedTo
}

// LogBlockedAccess 记录被阻止的访问尝试
func (v *WhitelistValidator) LogBlockedAccess(requestURL string) {
	fmt.Printf("访问被阻止 - URL: %s, 消息: %s\n", requestURL, v.config().BlockedMessage)
}

// GetAllowedDomains 获取允许访问的域名列表（用于调试）
func (v *WhitelistValidator) GetAllowedDomains() []string {
	return v.config().AllowedDomains
}

// UpdateConfig 更新白名单配置（运行时热更新）
// Deprecated
func (v *WhitelistValidator) UpdateConfig(newConfig func(...string) *config.WhitelistConfig) {
	v.config = newConfig
	fmt.Printf("白名单配置已更新，当前允许域名数量: %d\n", len(v.config().AllowedDomains))
}

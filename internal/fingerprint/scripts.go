// Package fingerprint 提供浏览器指纹伪装功能
// 包含静态脚本管理和动态脚本生成
package fingerprint

import (
	"embed"
)

// ScriptManager 指纹脚本管理器
type ScriptManager struct {
	resources      *embed.FS
	fingerprintJS  string
	websocketFixJS string
	headersFixJS   string
	isScriptLoaded bool
}

// NewScriptManager 创建新的脚本管理器实例
func NewScriptManager(resources *embed.FS) *ScriptManager {
	return &ScriptManager{
		resources:      resources,
		fingerprintJS:  "",
		websocketFixJS: "",
		headersFixJS:   "",
		isScriptLoaded: false,
	}
}

// LoadFingerprintScript 加载浏览器指纹伪装脚本
func (sm *ScriptManager) LoadFingerprintScript() error {
	// 从嵌入的资源中加载指纹伪装脚本
	if data, err := sm.resources.ReadFile("resources/fingerprint.js"); err == nil {
		sm.fingerprintJS = string(data)
		sm.isScriptLoaded = true
	} else {
		sm.fingerprintJS = ""
		sm.isScriptLoaded = false
	}

	// 加载WebSocket修复脚本
	if data, err := sm.resources.ReadFile("resources/websocket-fix.js"); err == nil {
		sm.websocketFixJS = string(data)
	} else {
		sm.websocketFixJS = ""
	}

	// 加载HTTP头部修复脚本
	if data, err := sm.resources.ReadFile("resources/headers-fix.js"); err == nil {
		sm.headersFixJS = string(data)
	} else {
		sm.headersFixJS = ""
	}

	return nil
}

// GetStaticScript 获取静态指纹脚本
func (sm *ScriptManager) GetStaticScript() string {
	return sm.fingerprintJS
}

// GetWebSocketFixScript 获取WebSocket修复脚本
func (sm *ScriptManager) GetWebSocketFixScript() string {
	return sm.websocketFixJS
}

// GetHeadersFixScript 获取HTTP头部修复脚本
func (sm *ScriptManager) GetHeadersFixScript() string {
	return sm.headersFixJS
}

// IsScriptLoaded 检查脚本是否已加载
func (sm *ScriptManager) IsScriptLoaded() bool {
	return sm.isScriptLoaded
}

// GetScriptInfo 获取脚本信息（用于调试）
func (sm *ScriptManager) GetScriptInfo() map[string]interface{} {
	return map[string]interface{}{
		"script_loaded": sm.isScriptLoaded,
		"script_length": len(sm.fingerprintJS),
		"has_content":   sm.fingerprintJS != "",
	}
}

// ReloadScript 重新加载脚本
func (sm *ScriptManager) ReloadScript() error {
	return sm.LoadFingerprintScript()
}

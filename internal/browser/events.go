// Package browser 浏览器事件处理
// 处理页面加载、导航、IPC通信等浏览器相关事件
package browser

import (
	"cef/internal/config"
	"cef/internal/fingerprint"
	"cef/internal/security"
	"fmt"
	"strings"
	"time"

	"github.com/energye/energy/v2/cef"
	"github.com/energye/energy/v2/cef/ipc"
	"github.com/energye/energy/v2/consts"
	"github.com/energye/golcl/lcl"
	"github.com/energye/golcl/lcl/rtl/version"
)

// 需要删除的HTTP头部
var needRemoveHeaderKey = []string{"DNT", "energy"}

// EventHandler 浏览器事件处理器
type EventHandler struct {
	browserConfig      *config.BrowserConfig
	whitelistValidator *security.WhitelistValidator
	scriptManager      *fingerprint.ScriptManager
	scriptGenerator    *fingerprint.Generator
	lastRedirectURL    string // 最后一次重定向的URL，用于防止循环
	redirectCount      int    // 重定向次数计数器
}

// NewEventHandler 创建新的事件处理器实例
func NewEventHandler(
	browserConfig *config.BrowserConfig,
	whitelistValidator *security.WhitelistValidator,
	scriptManager *fingerprint.ScriptManager,
	scriptGenerator *fingerprint.Generator,
) *EventHandler {
	return &EventHandler{
		browserConfig:      browserConfig,
		whitelistValidator: whitelistValidator,
		scriptManager:      scriptManager,
		scriptGenerator:    scriptGenerator,
	}
}

// SetupEvents 设置浏览器事件处理
func (h *EventHandler) SetupEvents(event *cef.BrowserEvent, window cef.IBrowserWindow) {
	// 监听来自前端JavaScript的"count"消息
	// 当前端调用ipc.emit("count", [count++])时，这个函数会被触发
	ipc.On("count", func(value int) {
		// 在控制台打印接收到的计数值
		println("count", value)
	})

	// 设置资源加载前的回调，用于修改请求头
	event.SetOnBeforeResourceLoad(func(sender lcl.IObject, browser *cef.ICefBrowser, frame *cef.ICefFrame, request *cef.ICefRequest, callback *cef.ICefCallback, result *consts.TCefReturnValue, window cef.IBrowserWindow) {
		// 获取并清理原有头部映射
		headerMap := request.GetHeaderMap()

		// 删除无用的头并重置header，同时去重
		cleanedHeaders := h.RemoveKey(headerMap, needRemoveHeaderKey)

		// 额外去重处理：确保每个头部键只出现一次
		deduplicatedHeaders := h.DeduplicateHeaders(cleanedHeaders)
		// 设置清理后的头部
		request.SetHeaderMap(deduplicatedHeaders)

		// 从UA中提取平台信息
		userAgent := h.browserConfig.Basic.UserAgent
		var platformValue string

		if strings.Contains(userAgent, "Windows") {
			platformValue = "Windows"
		} else if strings.Contains(userAgent, "Macintosh") {
			platformValue = "macOS"
		} else if strings.Contains(userAgent, "Linux") {
			platformValue = "Linux"
		} else if strings.Contains(userAgent, "Android") {
			platformValue = "Android"
		} else if strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad") {
			platformValue = "iOS"
		} else {
			platformValue = "Windows" // 默认平台
		}
		// 直接设置关键头部到request，确保生效
		request.SetHeaderByName("sec-ch-ua-platform", `"`+platformValue+`"`, true)
	})

	// 设置页面加载完成事件的处理函数
	// 当浏览器页面加载完成后会触发此事件
	event.SetOnLoadEnd(func(sender lcl.IObject, browser *cef.ICefBrowser, frame *cef.ICefFrame, httpStatusCode int32, window cef.IBrowserWindow) {
		h.handlePageLoad(browser, frame, httpStatusCode, window)
		if h.browserConfig.Proxy.Debug {
			window.Chromium().ExecuteJavaScript(`fetch('https://ifconfig.io/ip')
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.text();
    })
    .then(data => {
        console.log(data);
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });`, "", frame, 0)
		}
	})

	event.SetOnBeforeBrowser(func(sender lcl.IObject, browser *cef.ICefBrowser, frame *cef.ICefFrame, request *cef.ICefRequest, userGesture, isRedirect bool, window cef.IBrowserWindow) bool {
		requestContext := browser.GetRequestContext()
		if h.browserConfig.Proxy.Url != "" {
			proxyDict := cef.DictionaryValueRef.New()
			proxyDict.SetString("mode", h.browserConfig.Proxy.Mode)
			proxyDict.SetString("server", h.browserConfig.Proxy.Url)
			proxy := cef.ValueRef.New()
			proxy.SetDictionary(proxyDict)
			requestContext.SetPreference("proxy", proxy)
		}
		return false
	})

	window.Chromium().SetOnGetAuthCredentials(func(sender lcl.IObject, browser *cef.ICefBrowser, originUrl string, isProxy bool, host string, port int32, realm, scheme string, callback *cef.ICefAuthCallback) bool {
		if isProxy {
			callback.Cont(h.browserConfig.Proxy.Username, h.browserConfig.Proxy.Password)
			return true
		}
		return false
	})

	fmt.Println("✅ 浏览器事件处理器设置完成")
}

// RemoveKey 从StringMultiMap中删除指定的key
func (h *EventHandler) RemoveKey(header *cef.ICefStringMultiMap, keysToRemove []string) *cef.ICefStringMultiMap {
	// 临时存储保留的键值对
	preservedData := make(map[string][]string)

	// 遍历所有数据
	size := header.GetSize()
	var ketMapping = make(map[string]struct{}, len(keysToRemove))
	for _, s := range keysToRemove {
		ketMapping[s] = struct{}{}
	}
	for i := uint32(0); i < size; i++ {
		key := header.GetKey(i)
		value := header.GetValue(i)

		// 跳过要删除的键
		if _, ok := ketMapping[key]; ok {
			continue
		}

		// 保存其他键值对
		preservedData[key] = append(preservedData[key], value)
	}
	// 清空原数据
	header.Clear()
	// 重新添加保留的数据
	for key, values := range preservedData {
		for _, value := range values {
			header.Append(key, value)
		}
	}
	return header
}

// DeduplicateHeaders 去重StringMultiMap中的头部值
func (h *EventHandler) DeduplicateHeaders(header *cef.ICefStringMultiMap) *cef.ICefStringMultiMap {
	// 临时存储保留的键值对
	preservedData := make(map[string][]string)

	// 遍历所有数据
	size := header.GetSize()
	for i := uint32(0); i < size; i++ {
		key := header.GetKey(i)
		value := header.GetValue(i)

		// 保存其他键值对
		preservedData[key] = append(preservedData[key], value)
	}
	// 清空原数据
	header.Clear()
	// 重新添加保留的数据
	for key, values := range preservedData {
		// 确保每个键只保留一个值
		if len(values) > 0 {
			header.Append(key, values[0])
		}
	}
	return header
}

// handlePageLoad 处理页面加载完成事件
func (h *EventHandler) handlePageLoad(browser *cef.ICefBrowser, frame *cef.ICefFrame, httpStatusCode int32, window cef.IBrowserWindow) {
	currentURL := frame.Url()

	// 输出调试信息帮助定位问题
	fmt.Printf("页面加载事件 - URL: %s, 状态码: %d\n", currentURL, httpStatusCode)

	// 检查URL是否被允许访问（优先检查，避免不必要的脚本注入）
	if currentURL != "" && currentURL != "about:blank" && !h.whitelistValidator.IsURLAllowed(currentURL) {
		h.handleBlockedURL(browser, currentURL)
		return
	}

	// 仅对允许的URL进行指纹注入
	fmt.Printf("开始注入指纹脚本 - URL: %s\n", currentURL)
	h.injectFingerprintScripts(browser)

	// 延迟补强注入（仅一次）
	go func() {
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("延迟补强注入 - URL: %s\n", currentURL)
		h.injectFingerprintScripts(browser)
	}()

	// 发送系统信息到前端
	h.sendSystemInfo(window)
}

// handleBlockedURL 处理被阻止的URL访问
func (h *EventHandler) handleBlockedURL(browser *cef.ICefBrowser, currentURL string) {
	// 防止重定向循环：检查是否与上次重定向目标相同
	redirectURL := h.whitelistValidator.GetRedirectURL()
	if currentURL == h.lastRedirectURL || h.redirectCount > 3 {
		// 避免无限重定向循环
		return
	}

	fmt.Printf("检测到不允许的URL: %s，准备重定向\n", currentURL)
	h.whitelistValidator.LogBlockedAccess(currentURL)

	if redirectURL != "" && redirectURL != currentURL {
		h.lastRedirectURL = currentURL
		h.redirectCount++

		fmt.Printf("重定向到: %s\n", redirectURL)
		browser.MainFrame().LoadUrl(redirectURL)

		// 重置计数器（延迟重置）
		go func() {
			time.Sleep(5 * time.Second)
			h.redirectCount = 0
		}()
	}
}

// injectFingerprintScripts 注入指纹伪装脚本
func (h *EventHandler) injectFingerprintScripts(browser *cef.ICefBrowser) {
	// 注入HTTP头部修复脚本
	headersFixScript := h.scriptManager.GetHeadersFixScript()
	if headersFixScript != "" {
		browser.MainFrame().ExecuteJavaScript(headersFixScript, "", 0)
		fmt.Printf("HTTP头部伪装已启用\n")
	}

	// 注入WebSocket修复脚本
	websocketFixScript := h.scriptManager.GetWebSocketFixScript()
	if websocketFixScript != "" {
		browser.MainFrame().ExecuteJavaScript(websocketFixScript, "", 0)
		fmt.Printf("WebSocket优雅错误处理已启用\n")
	}

	// 注入CORS禁用脚本（在指纹脚本之前）
	corsScript := `
		console.log('开始设置 CORS 禁用...');
		
		// 禁用 Fetch CORS 检查
		if (window.fetch) {
			const originalFetch = window.fetch;
			window.fetch = function(url, options = {}) {
				options.mode = 'cors';
				options.credentials = 'include';
				return originalFetch(url, options).catch(error => {
					console.log('Fetch CORS 错误已被忽略:', error);
					return new Response('{}', { status: 200, statusText: 'OK' });
				});
			};
		}
		
		// WebSocket 连接增强处理
		if (window.WebSocket) {
			const OriginalWebSocket = window.WebSocket;
			window.WebSocket = function(url, protocols) {
				console.log('WebSocket 连接请求:', url);
				
				// 创建增强的WebSocket实例
				const ws = new OriginalWebSocket(url, protocols);
				
				// 增强错误处理
				const originalOnError = ws.onerror;
				ws.onerror = function(event) {
					console.warn('WebSocket 连接失败，尝试恢复...', event);
					if (originalOnError) originalOnError.call(this, event);
				};
				
				// 成功连接日志
				const originalOnOpen = ws.onopen;
				ws.onopen = function(event) {
					console.log('WebSocket 连接成功:', url);
					if (originalOnOpen) originalOnOpen.call(this, event);
				};
				
				return ws;
			};
			
			// 保持原型和常量
			window.WebSocket.prototype = OriginalWebSocket.prototype;
			window.WebSocket.CONNECTING = OriginalWebSocket.CONNECTING;
			window.WebSocket.OPEN = OriginalWebSocket.OPEN;
			window.WebSocket.CLOSING = OriginalWebSocket.CLOSING;
			window.WebSocket.CLOSED = OriginalWebSocket.CLOSED;
		}
		
		console.log('CORS 禁用和 WebSocket 增强设置完成');
	`
	browser.MainFrame().ExecuteJavaScript(corsScript, "", 0)

	// 注入静态指纹脚本
	if h.scriptManager.IsScriptLoaded() {
		staticScript := h.scriptManager.GetStaticScript()
		browser.MainFrame().ExecuteJavaScript(staticScript, "", 0)
	}

	// 注入动态基础指纹脚本
	basicScript := h.scriptGenerator.GenerateBasicScript()
	browser.MainFrame().ExecuteJavaScript(basicScript, "", 0)

	// 注入高级指纹脚本
	advancedScript := h.scriptGenerator.GenerateAdvancedScript()
	browser.MainFrame().ExecuteJavaScript(advancedScript, "", 0)

	fmt.Printf("指纹伪装、CORS禁用、WebSocket优雅处理和HTTP头部伪装已应用\n")
}

// sendSystemInfo 发送系统信息到前端
func (h *EventHandler) sendSystemInfo(window cef.IBrowserWindow) {
	// 获取操作系统版本信息
	osVersion := version.OSVersion.ToString()
	// 静默执行，不输出日志
	// println("osInfo", osVersion)

	// 通过IPC将操作系统信息发送给前端JavaScript
	// 前端可以通过ipc.on("osInfo", function(os){...})接收这个信息
	ipc.Emit("osInfo", osVersion)

	// 判断窗口类型并设置相应的字符串标识
	var windowType string
	if window.IsLCL() {
		// LCL类型窗口（Lazarus Component Library）
		windowType = "LCL"
	} else {
		// VF类型窗口（可能是ViewFrame）
		windowType = "VF"
	}

	// 通过IPC将窗口类型信息发送给前端JavaScript
	// 前端可以通过ipc.on("windowType", function(type){...})接收这个信息
	ipc.Emit("windowType", windowType)
}

// UpdateConfigs 更新配置（运行时热更新）
func (h *EventHandler) UpdateConfigs(
	browserConfig *config.BrowserConfig,
	whitelistValidator *security.WhitelistValidator,
) {
	h.browserConfig = browserConfig
	h.whitelistValidator = whitelistValidator
	h.scriptGenerator.UpdateConfig(browserConfig)

	fmt.Println("事件处理器配置已更新")
}

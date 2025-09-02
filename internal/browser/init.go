// Package browser 浏览器初始化
// 负责CEF浏览器的初始化、配置和启动
package browser

import (
	"cef/internal/config"
	"embed"
	"fmt"

	"github.com/energye/energy/v2/cef"
	"github.com/energye/energy/v2/pkgs/assetserve"
)

// Initializer 浏览器初始化器
type Initializer struct {
	resources     *embed.FS
	browserConfig *config.BrowserConfig
	eventHandler  *EventHandler
}

// NewInitializer 创建新的浏览器初始化器实例
func NewInitializer(
	resources *embed.FS,
	browserConfig *config.BrowserConfig,
	eventHandler *EventHandler,
) *Initializer {
	return &Initializer{
		resources:     resources,
		browserConfig: browserConfig,
		eventHandler:  eventHandler,
	}
}

// Initialize 初始化CEF浏览器
func (init *Initializer) Initialize() *cef.TCEFApplication {
	// 创建一个Energy应用程序实例
	app := cef.NewApplication()

	// 设置全局User-Agent（影响所有HTTP请求）
	app.SetUserAgent(init.browserConfig.Basic.UserAgent)
	fmt.Printf("已设置HTTP User-Agent: %s\n", init.browserConfig.Basic.UserAgent)

	// 设置平台信息
	app.AddCustomCommandLine("--user-agent", init.browserConfig.Basic.UserAgent)

	// 设置接受语言
	app.AddCustomCommandLine("--lang", "zh-CN")
	app.AddCustomCommandLine("--accept-lang", init.browserConfig.Basic.AcceptLanguage)

	// 设置CEF命令行开关来解决CORS和安全问题
	app.AddCustomCommandLine("--disable-web-security", "")
	app.AddCustomCommandLine("--allow-running-insecure-content", "")
	app.AddCustomCommandLine("--ignore-certificate-errors", "")
	app.AddCustomCommandLine("--ignore-ssl-errors", "")
	app.AddCustomCommandLine("--ignore-urlfetcher-cert-requests", "")
	app.AddCustomCommandLine("--disable-extensions", "")
	app.AddCustomCommandLine("--disable-plugins", "")
	app.AddCustomCommandLine("--disable-default-apps", "")
	app.AddCustomCommandLine("--disable-background-timer-throttling", "")

	// WebSocket深度支持配置 - 解决状态码200问题
	app.AddCustomCommandLine("--enable-websockets", "")
	app.AddCustomCommandLine("--enable-experimental-web-platform-features", "")
	app.AddCustomCommandLine("--disable-site-isolation-trials", "")
	app.AddCustomCommandLine("--allow-websocket-upgrade-on-any-port", "")
	app.AddCustomCommandLine("--disable-background-networking", "")
	app.AddCustomCommandLine("--disable-sync", "")
	app.AddCustomCommandLine("--disable-translate", "")
	app.AddCustomCommandLine("--no-proxy-server", "")
	app.AddCustomCommandLine("--disable-renderer-backgrounding", "")

	// 统一的禁用功能配置（避免重复）
	app.AddCustomCommandLine("--disable-features", "VizDisplayCompositor,SiteIsolation,TranslateUI,BackgroundSync")
	fmt.Println("已配置CEF安全选项以避免CORS错误和WebSocket连接问题")

	// 配置浏览器窗口
	init.configureBrowserWindow()

	// 配置静态资源服务器
	init.configureAssetServer()

	// 设置浏览器窗口初始化时的回调函数
	// 用于配置IPC通信和各种事件处理
	cef.BrowserWindow.SetBrowserInit(func(event *cef.BrowserEvent, window cef.IBrowserWindow) {
		init.eventHandler.SetupEvents(event, window)
	})

	fmt.Println("CEF初始化完成")
	return app
}

// configureBrowserWindow 配置浏览器窗口
func (init *Initializer) configureBrowserWindow() {
	// 设置浏览器窗口默认加载的URL（从配置文件读取）
	cef.BrowserWindow.Config.Url = init.browserConfig.App.DefaultURL

	// 设置窗口标题
	cef.BrowserWindow.Config.Title = init.browserConfig.App.WindowTitle

	// 设置窗口大小
	cef.BrowserWindow.Config.Width = int32(init.browserConfig.Screen.Width)
	cef.BrowserWindow.Config.Height = int32(init.browserConfig.Screen.Height)

	// 注意：Energy框架的UserAgent和其他浏览器指纹配置
	// 需要通过JavaScript注入的方式实现，而不是在这里配置
}

// configureAssetServer 配置内置静态资源服务器
func (init *Initializer) configureAssetServer() {
	// 配置内置静态资源服务器的安全验证头
	// 这是一种简单的安全机制，防止未授权访问静态资源
	assetserve.AssetsServerHeaderKeyName = "energy"
	assetserve.AssetsServerHeaderKeyValue = "energy"

	// 设置浏览器进程启动完成后的回调函数
	// 当浏览器进程启动完成后，会执行这个回调来启动内置的HTTP服务器
	cef.SetBrowserProcessStartAfterCallback(func(b bool) {
		// 创建一个新的静态资源HTTP服务器实例
		server := assetserve.NewAssetsHttpServer()
		// 设置HTTP服务器监听端口为22022
		server.PORT = 22022
		// 指定资源文件夹名称，与embed指令中的目录名对应
		server.AssetsFSName = "resources"
		// 将嵌入的文件系统赋值给服务器，使其能够提供静态资源服务
		server.Assets = init.resources
		// 在新的goroutine中启动HTTP服务器（非阻塞方式）
		go server.StartHttpServer()

		// 输出调试信息
		fmt.Printf("静态资源HTTP服务器已启动在端口22022\n")
		fmt.Printf("访问 URL: http://localhost:22022/index.html\n")
	})
}

// UpdateConfig 更新浏览器配置
func (init *Initializer) UpdateConfig(newConfig *config.BrowserConfig) {
	init.browserConfig = newConfig
	fmt.Println("浏览器初始化器配置已更新")
}

// GetConfigSummary 获取当前配置摘要
func (init *Initializer) GetConfigSummary() map[string]interface{} {
	return map[string]interface{}{
		"default_url":   init.browserConfig.App.DefaultURL,
		"window_title":  init.browserConfig.App.WindowTitle,
		"screen_width":  init.browserConfig.Screen.Width,
		"screen_height": init.browserConfig.Screen.Height,
		"user_agent":    init.browserConfig.Basic.UserAgent,
		"language":      init.browserConfig.Basic.AcceptLanguage,
	}
}

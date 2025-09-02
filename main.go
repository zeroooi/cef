// 基于Energy框架和CEF的受控浏览器客户端
// 功能：模拟真实浏览器环境，限制访问指定网站，提供浏览器指纹伪装
package main

import (
	"embed" // Go内置的文件嵌入功能
	"log"   // 日志记录

	"cef/internal/browser"     // 浏览器初始化和事件处理
	"cef/internal/config"      // 配置管理
	"cef/internal/fingerprint" // 指纹伪装
	"cef/internal/security"    // 安全控制（白名单等）

	"github.com/energye/energy/v2/cef" // Energy CEF核心包
)

// 使用Go的embed指令将resources目录下的所有文件嵌入到程序中
// 编译后的可执行文件将包含所有静态资源（HTML、CSS、JS、图片等）
//
//go:embed resources
var resources embed.FS

// 应用程序主入口函数
func main() {
	// 1. 加载配置文件
	configLoader := config.NewLoader()
	if err := configLoader.LoadAll(); err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}
	log.Println("配置加载成功")

	// 获取配置实例
	browserConfig := configLoader.GetBrowserConfig()
	whitelistConfig := configLoader.GetWhitelistConfig()

	// 2. 初始化安全控制模块
	whitelistValidator := security.NewWhitelistValidator(whitelistConfig)
	log.Println("安全控制模块初始化完成")

	// 3. 初始化指纹伪装模块
	scriptManager := fingerprint.NewScriptManager(&resources)
	if err := scriptManager.LoadFingerprintScript(); err != nil {
		log.Printf(" 警告：静态指纹脚本加载失败，将仅使用动态脚本: %v", err)
	} else {
		log.Println("指纹伪装脚本加载成功")
	}

	scriptGenerator := fingerprint.NewGenerator(browserConfig)
	log.Println("指纹伪装模块初始化完成")

	// 4. 初始化浏览器事件处理器
	eventHandler := browser.NewEventHandler(
		browserConfig,
		whitelistValidator,
		scriptManager,
		scriptGenerator,
	)
	log.Println("浏览器事件处理器初始化完成")

	// 5. 初始化浏览器
	browserInit := browser.NewInitializer(&resources, browserConfig, eventHandler)

	log.Println("正在初始化 CEF 浏览器...")
	// 初始化CEF框架（只能调用一次）
	cef.GlobalInit(nil, &resources)

	app := browserInit.Initialize()

	log.Println("启动 CEF 应用...")
	// 6. 启动并运行应用程序
	// 这会阻塞主线程直到应用程序退出
	cef.Run(app)
}

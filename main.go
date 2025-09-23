// 基于Energy框架和CEF的受控浏览器客户端
// 功能：模拟真实浏览器环境，限制访问指定网站，提供浏览器指纹伪装
package main

import (
	"cef/internal/browser"     // 浏览器初始化和事件处理
	"cef/internal/config"      // 配置管理
	"cef/internal/fingerprint" // 指纹伪装
	"cef/internal/security"    // 安全控制（白名单等）
	"cef/pkg/external/aegis"
	"embed"                            // Go内置的文件嵌入功能
	"github.com/energye/energy/v2/cef" // Energy CEF核心包
	"log"                              // 日志记录
)

// 使用Go的embed指令将resources目录下的所有文件嵌入到程序中
// 编译后的可执行文件将包含所有静态资源（HTML、CSS、JS、图片等）
//
//go:embed resources
var resources embed.FS

//go:embed config
var cfg embed.FS

// 应用程序主入口函数
func main() {
	// 1. 加载配置文件
	configLoader := config.NewLoader(&cfg)
	if err := configLoader.LoadAll(); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}
	log.Println("配置加载成功")

	aegis.SetDefault(aegis.NewAegisClient(configLoader.ExternalConfig.AegisAddr.Mode))

	// 获取配置实例
	browserConfigLoader := configLoader.GetBrowserConfigLoader()
	whitelistConfigLoader := configLoader.GetWhitelistConfigLoader()

	// 2. 初始化安全控制模块
	whitelistValidator := security.NewWhitelistValidator(whitelistConfigLoader)
	log.Println("安全控制模块初始化完成")

	// 3. 初始化指纹伪装模块
	scriptManager := fingerprint.NewScriptManager(&resources)
	if err := scriptManager.LoadFingerprintScript(); err != nil {
		log.Printf(" 警告：静态指纹脚本加载失败，将仅使用动态脚本: %v", err)
	} else {
		log.Println("指纹伪装脚本加载成功")
	}

	scriptGenerator := fingerprint.NewGenerator(browserConfigLoader)
	log.Println("指纹伪装模块初始化完成")

	notifyAccountChangeChan := make(chan string, 1)
	// 4. 初始化浏览器事件处理器
	eventHandler := browser.NewEventHandler(
		browserConfigLoader,
		whitelistValidator,
		scriptManager,
		scriptGenerator,
		notifyAccountChangeChan,
	)
	defer eventHandler.Close()

	log.Println("浏览器事件处理器初始化完成")

	// 5. 初始化浏览器
	browserInit := browser.NewInitializer(&resources, browserConfigLoader(), eventHandler)

	log.Println("正在初始化 CEF 浏览器...")
	// 初始化CEF框架（只能调用一次）
	cef.GlobalInit(nil, &resources)
	cef.BrowserWindow.Config.IconFS = "resources/icon.png"
	app := browserInit.Initialize()

	//go func() {
	//	for account := range notifyAccountChangeChan {
	//		fmt.Println("收到AccountChange信息， account:", account)
	//		browserConfig := browserConfigLoader(account)
	//		fmt.Println("开始设置app.SetUserAgent:", browserConfig.Basic.UserAgent)
	//		// 设置全局User-Agent（影响所有HTTP请求）
	//		app.SetUserAgent(browserConfig.Basic.UserAgent)
	//
	//		// 设置平台信息
	//		//app.AddCustomCommandLine("--user-agent", browserConfig.Basic.UserAgent)
	//		fmt.Println("成功设置app.SetUserAgent")
	//	}
	//	fmt.Println("notifyAccountChangeChan关闭")
	//}()

	log.Println("启动 CEF 应用...")
	// 6. 启动并运行应用程序
	// 这会阻塞主线程直到应用程序退出
	cef.Run(app)
}

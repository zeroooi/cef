package main

import (
	"cef/internal/config"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func main() {
	fmt.Println("=== 配置加载测试 ===")

	// 直接使用Viper测试
	v := viper.New()
	v.SetConfigName("browser_config")
	v.SetConfigType("json")
	v.AddConfigPath("../config")

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Viper读取配置失败: %v\n", err)
		return
	}

	fmt.Println("✅ Viper成功读取配置文件")

	// 检查Viper读取的原始值
	fmt.Printf("Viper - user_agent: '%s'\n", v.GetString("basic.user_agent"))
	fmt.Printf("Viper - accept_language: '%s'\n", v.GetString("basic.accept_language"))
	fmt.Printf("Viper - default_url: '%s'\n", v.GetString("app.default_url"))
	fmt.Printf("Viper - window_title: '%s'\n", v.GetString("app.window_title"))
	fmt.Printf("Viper - cpu_cores: %d\n", v.GetInt("hardware.cpu_cores"))
	fmt.Printf("Viper - device_memory: %d\n", v.GetInt("hardware.device_memory"))

	// 测试直接Unmarshal
	var testConfig config.BrowserConfig
	if err := v.Unmarshal(&testConfig); err != nil {
		fmt.Printf("直接Unmarshal失败: %v\n", err)
	} else {
		fmt.Println("✅ 直接Unmarshal成功")
		fmt.Printf("Unmarshal - user_agent: '%s'\n", testConfig.Basic.UserAgent)
		fmt.Printf("Unmarshal - default_url: '%s'\n", testConfig.App.DefaultURL)
	}

	// 使用配置加载器测试
	fmt.Println("\n=== 使用配置加载器 ===")

	// 切换到项目根目录，因为配置加载器需要在项目根目录下运行
	originalDir, _ := os.Getwd()
	os.Chdir("..")
	defer os.Chdir(originalDir)

	loader := config.NewLoader()
	if err := loader.LoadBrowserConfig(); err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		return
	}

	browserConfig := loader.GetBrowserConfig()
	fmt.Printf("Loader - user_agent: '%s'\n", browserConfig.Basic.UserAgent)
	fmt.Printf("Loader - default_url: '%s'\n", browserConfig.App.DefaultURL)
}

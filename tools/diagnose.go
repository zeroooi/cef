// CEF 环境诊断工具
// 检查 macOS 上 CEF 应用运行所需的环境和权限
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	fmt.Println("=== CEF 环境诊断工具 ===")
	fmt.Printf("操作系统: %s\n", runtime.GOOS)
	fmt.Printf("架构: %s\n", runtime.GOARCH)
	fmt.Printf("Go 版本: %s\n", runtime.Version())
	fmt.Println()

	// 检查当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ 无法获取当前目录: %v\n", err)
	} else {
		fmt.Printf("✅ 当前工作目录: %s\n", wd)
	}

	// 检查配置文件
	checkFile("✅ 配置文件", "../config/browser_config.json")
	checkFile("✅ 白名单配置", "../config/whitelist.json")
	checkFile("✅ 指纹脚本", "../resources/fingerprint.js")
	checkFile("✅ 默认页面", "../resources/index.html")

	// 检查进程
	fmt.Println("\n=== 进程检查 ===")
	checkProcess("CEF 相关进程", "cef")
	checkProcess("Energy 进程", "energy")
	checkProcess("Chromium 进程", "Chromium")

	// 检查网络端口
	fmt.Println("\n=== 网络检查 ===")
	checkPort("HTTP 服务端口", "22022")

	// 检查权限（macOS特定）
	fmt.Println("\n=== 权限检查 ===")
	if runtime.GOOS == "darwin" {
		checkMacOSPermissions()
	}

	fmt.Println("\n=== 诊断完成 ===")
	fmt.Println("如果遇到权限问题，请：")
	fmt.Println("1. 运行清理脚本: ../cleanup.sh")
	fmt.Println("2. 重新启动应用: cd .. && go run main.go env=dev")
	fmt.Println("3. 在弹出的权限对话框中点击'允许'")
}

func checkFile(desc, path string) {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("✅ %s: %s\n", desc, path)
	} else {
		fmt.Printf("❌ %s: %s (不存在)\n", desc, path)
	}
}

func checkProcess(desc, name string) {
	cmd := exec.Command("pgrep", "-f", name)
	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		fmt.Printf("✅ %s: 无运行进程\n", desc)
	} else {
		fmt.Printf("⚠️  %s: 发现运行进程\n", desc)
		fmt.Printf("   PID: %s", string(output))
	}
}

func checkPort(desc, port string) {
	cmd := exec.Command("lsof", "-i", ":"+port)
	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		fmt.Printf("✅ %s (%s): 端口空闲\n", desc, port)
	} else {
		fmt.Printf("⚠️  %s (%s): 端口被占用\n", desc, port)
	}
}

func checkMacOSPermissions() {
	// 检查是否在应用程序文件夹中
	wd, _ := os.Getwd()
	if wd == "/Applications" || wd == "/System/Applications" {
		fmt.Println("⚠️  当前在系统应用目录，可能需要管理员权限")
	} else {
		fmt.Println("✅ 当前目录权限正常")
	}

	// 检查 Gatekeeper 状态
	cmd := exec.Command("spctl", "--status")
	output, err := cmd.Output()
	if err == nil {
		fmt.Printf("ℹ️  Gatekeeper 状态: %s", string(output))
	}

	fmt.Println("💡 macOS 权限建议：")
	fmt.Println("   - 在'系统偏好设置 > 安全性与隐私'中允许应用运行")
	fmt.Println("   - 如果提示需要密码，请输入您的用户密码")
	fmt.Println("   - 确保应用有访问网络的权限")
}

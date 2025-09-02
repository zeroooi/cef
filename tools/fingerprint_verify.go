// 指纹模拟效果验证工具
// 用于快速检测指纹模拟是否成功
package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	fmt.Println("🔍 指纹模拟效果验证工具")
	fmt.Println("========================")

	// 检查操作系统
	if runtime.GOOS != "darwin" {
		fmt.Println("⚠️  此工具主要用于macOS平台的指纹模拟验证")
	}

	fmt.Println("📋 验证步骤：")
	fmt.Println("1. 确保CEF应用正在运行")
	fmt.Println("2. 在应用中访问指纹测试网站")
	fmt.Println("3. 对比实际结果与预期配置")

	fmt.Println("\n🌐 推荐的在线测试网站：")
	sites := []struct {
		name string
		url  string
		desc string
	}{
		{"AmIUnique", "https://amiunique.org/", "综合指纹唯一性检测"},
		{"BrowserLeaks", "https://browserleaks.com/", "全面浏览器泄露检测"},
		{"Canvas测试", "https://browserleaks.com/canvas", "Canvas指纹专项测试"},
		{"WebGL测试", "https://webglreport.com/", "WebGL信息检测"},
		{"本地测试页", "http://localhost:22022/fingerprint-test.html", "自建指纹测试页面"},
	}

	for _, site := range sites {
		fmt.Printf("  • %s: %s\n    %s\n", site.name, site.url, site.desc)
	}

	fmt.Println("\n🎯 关键验证点：")
	checks := []struct {
		item     string
		expected string
		risk     string
	}{
		{"User-Agent", "Windows Chrome 120", "显示macOS会暴露真实系统"},
		{"平台信息", "Win32", "显示Darwin会暴露macOS"},
		{"屏幕分辨率", "1920x1080", "显示Mac分辨率会暴露设备"},
		{"CPU核心", "8核", "显示真实Mac CPU配置"},
		{"WebGL渲染器", "Intel集显", "显示真实Mac显卡信息"},
		{"Canvas指纹", "每次略有不同", "完全一致表明未注入噪声"},
	}

	for _, check := range checks {
		fmt.Printf("  ✅ %s: %s\n     ❌ %s\n", check.item, check.expected, check.risk)
	}

	fmt.Println("\n📊 快速验证脚本：")
	fmt.Println("如果应用正在运行，可以尝试在浏览器控制台执行：")
	fmt.Println(`
console.log('=== 快速指纹检查 ===');
console.log('User-Agent:', navigator.userAgent);
console.log('平台:', navigator.platform);
console.log('语言:', navigator.language);
console.log('屏幕:', screen.width + 'x' + screen.height);
console.log('CPU核心:', navigator.hardwareConcurrency);
console.log('设备内存:', navigator.deviceMemory + 'GB');

// WebGL检测
const canvas = document.createElement('canvas');
const gl = canvas.getContext('webgl');
if (gl) {
    console.log('WebGL渲染器:', gl.getParameter(gl.RENDERER));
    console.log('WebGL供应商:', gl.getParameter(gl.VENDOR));
}
`)

	fmt.Println("\n⚡ 自动化测试：")
	fmt.Println("正在检查CEF进程...")

	// 检查CEF进程
	cmd := exec.Command("pgrep", "-f", "cef")
	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		fmt.Println("❌ 未检测到CEF进程运行")
		fmt.Println("   请先启动应用: cd .. && go run main.go env=dev")
	} else {
		fmt.Println("✅ 检测到CEF进程正在运行")

		// 等待一下让用户看到消息
		time.Sleep(2 * time.Second)

		// 尝试打开本地测试页面
		fmt.Println("🚀 尝试打开本地测试页面...")
		openCmd := exec.Command("open", "http://localhost:22022/fingerprint-test.html")
		if err := openCmd.Run(); err != nil {
			fmt.Println("⚠️  无法自动打开页面，请手动访问：")
			fmt.Println("   http://localhost:22022/fingerprint-test.html")
		} else {
			fmt.Println("✅ 已在默认浏览器中打开测试页面")
			fmt.Println("   请切换到CEF应用中查看指纹测试结果")
		}
	}

	fmt.Println("\n🛟 问题排查：")
	fmt.Println("如果发现指纹泄露，请检查：")
	fmt.Println("  1. 指纹脚本是否正确加载")
	fmt.Println("  2. 配置文件参数是否正确")
	fmt.Println("  3. JavaScript注入是否成功")
	fmt.Println("  4. 使用开发者工具查看Console错误")

	fmt.Println("\n📝 建议测试流程：")
	fmt.Println("  1. 运行: cd tools && go run fingerprint_verify.go")
	fmt.Println("  2. 在CEF应用中访问推荐的测试网站")
	fmt.Println("  3. 截图保存测试结果")
	fmt.Println("  4. 与配置参数进行对比验证")
}

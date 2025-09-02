// Package config 配置加载器实现
// 使用Viper组件实现配置文件的加载和解析
package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Loader 配置加载器
type Loader struct {
	browserConfig   *BrowserConfig
	whitelistConfig *WhitelistConfig
}

// NewLoader 创建新的配置加载器实例
func NewLoader() *Loader {
	return &Loader{
		browserConfig:   &BrowserConfig{},
		whitelistConfig: &WhitelistConfig{},
	}
}

// LoadAll 加载所有配置文件
func (l *Loader) LoadAll() error {
	// 加载浏览器配置
	if err := l.LoadBrowserConfig(); err != nil {
		return fmt.Errorf("加载浏览器配置失败: %v", err)
	}

	// 加载白名单配置
	if err := l.LoadWhitelistConfig(); err != nil {
		return fmt.Errorf("加载白名单配置失败: %v", err)
	}

	fmt.Println("配置加载完成")
	fmt.Printf("允许访问的域名: %v\n", l.whitelistConfig.AllowedDomains)

	return nil
}

// LoadBrowserConfig 使用Viper加载浏览器配置
func (l *Loader) LoadBrowserConfig() error {
	// 创建Viper实例用于浏览器配置
	v := viper.New()
	v.SetConfigName("browser_config") // 配置文件名（不包含扩展名）
	v.SetConfigType("json")           // 配置文件类型
	v.AddConfigPath("./config")       // 配置文件路径

	// 先读取配置文件，再设置默认值
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// 配置文件不存在，使用默认值
			fmt.Printf("浏览器配置文件不存在，使用默认配置: %v\n", err)
			l.setDefaultBrowserConfig(v)
		}
	} else {
		fmt.Println("成功加载浏览器配置文件")
		// 文件加载成功，不设置默认值以避免覆盖
	}

	// 将配置值映射到结构体
	// 由于Unmarshal有问题，直接使用手动映射方式
	fmt.Println("使用手动配置映射方式")
	l.setFallbackBrowserConfig(v)

	fmt.Printf("浏览器配置加载完成: User-Agent=%s, 默认URL=%s\n",
		l.browserConfig.Basic.UserAgent, l.browserConfig.App.DefaultURL)
	fmt.Printf("Canvas噪声: %v, WebGL渲染器: %s\n",
		l.browserConfig.Canvas.EnableNoise, l.browserConfig.WebGL.Renderer)

	return nil
}

// LoadWhitelistConfig 使用Viper加载白名单配置
func (l *Loader) LoadWhitelistConfig() error {
	// 创建Viper实例用于白名单配置
	v := viper.New()
	v.SetConfigName("whitelist") // 配置文件名（不包含扩展名）
	v.SetConfigType("json")      // 配置文件类型
	v.AddConfigPath("./config")  // 配置文件路径

	// 设置默认值
	v.SetDefault("allowed_domains", []string{"google.com", "www.google.com", "accounts.google.com"})
	v.SetDefault("blocked_message", "访问被限制：该网站不在允许访问列表中")
	v.SetDefault("redirect_blocked_to", "https://www.google.com")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// 配置文件不存在，使用默认值
			fmt.Printf("白名单配置文件不存在，使用默认配置: %v\n", err)
		}
	} else {
		fmt.Println("成功加载白名单配置文件")
	}

	// 将配置值映射到结构体
	l.whitelistConfig.AllowedDomains = v.GetStringSlice("allowed_domains")
	l.whitelistConfig.BlockedMessage = v.GetString("blocked_message")
	l.whitelistConfig.RedirectBlockedTo = v.GetString("redirect_blocked_to")

	fmt.Printf("白名单配置加载完成: 允许域名数量=%d\n", len(l.whitelistConfig.AllowedDomains))

	return nil
}

// GetBrowserConfig 获取浏览器配置
func (l *Loader) GetBrowserConfig() *BrowserConfig {
	return l.browserConfig
}

// GetWhitelistConfig 获取白名单配置
func (l *Loader) GetWhitelistConfig() *WhitelistConfig {
	return l.whitelistConfig
}

// setDefaultBrowserConfig 设置浏览器配置的默认值
func (l *Loader) setDefaultBrowserConfig(v *viper.Viper) {
	v.SetDefault("basic.user_agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	v.SetDefault("basic.accept_language", "zh-CN,zh;q=0.9,en;q=0.8")
	v.SetDefault("basic.timezone", "Asia/Shanghai")
	v.SetDefault("basic.platform", "Win32")
	v.SetDefault("basic.vendor", "Google Inc.")
	v.SetDefault("basic.product", "Gecko")

	v.SetDefault("screen.width", 1920)
	v.SetDefault("screen.height", 1080)
	v.SetDefault("screen.avail_width", 1920)
	v.SetDefault("screen.avail_height", 1040)
	v.SetDefault("screen.color_depth", 24)
	v.SetDefault("screen.pixel_depth", 24)
	v.SetDefault("screen.device_pixel_ratio", 1.0)

	v.SetDefault("hardware.cpu_cores", 8)
	v.SetDefault("hardware.device_memory", 8)
	v.SetDefault("hardware.max_touch_points", 0)
	v.SetDefault("hardware.vendor_sub", "")
	v.SetDefault("hardware.product_sub", "20030107")

	v.SetDefault("canvas.enable_noise", true)
	v.SetDefault("canvas.noise_level", 0.1)
	v.SetDefault("canvas.block_toDataURL", false)

	v.SetDefault("webgl.vendor", "Google Inc. (Intel)")
	v.SetDefault("webgl.renderer", "ANGLE (Intel, Intel(R) UHD Graphics 630 Direct3D11 vs_5_0 ps_5_0, D3D11)")
	v.SetDefault("webgl.version", "WebGL 1.0 (OpenGL ES 2.0 Chromium)")
	v.SetDefault("webgl.shading_language_version", "WebGL GLSL ES 1.0 (OpenGL ES GLSL ES 1.0 Chromium)")

	v.SetDefault("audio.enable_noise", true)
	v.SetDefault("audio.noise_level", 0.0001)

	v.SetDefault("webrtc.block_local_ip_leak", true)
	v.SetDefault("webrtc.fake_public_ip", "8.8.8.8")
	v.SetDefault("webrtc.block_data_channels", false)

	v.SetDefault("app.default_url", "https://www.google.com")
	v.SetDefault("app.window_title", "安全浏览器")

	// HTTP头部默认值
	v.SetDefault("headers.sec_ch_ua", `"Not:A=Brand";v="99", "Google Chrome";v="139", "Chromium";v="139"`)
	v.SetDefault("headers.sec_ch_ua_mobile", "?0")
	v.SetDefault("headers.sec_ch_ua_platform", `"macOS"`)
	v.SetDefault("headers.sec_ch_ua_full_version_list", `"Not:A=Brand";v="99.0.0.0", "Google Chrome";v="139.0.0.0", "Chromium";v="139.0.0.0"`)
	v.SetDefault("headers.sec_ch_ua_arch", `"x86"`)
	v.SetDefault("headers.sec_ch_ua_bitness", `"64"`)
	v.SetDefault("headers.sec_fetch_dest", "empty")
	v.SetDefault("headers.sec_fetch_mode", "cors")
	v.SetDefault("headers.sec_fetch_site", "same-origin")
	v.SetDefault("headers.cache_control", "no-cache")
	v.SetDefault("headers.pragma", "no-cache")
	v.SetDefault("headers.x_sw_cache", "7")
}

// setFallbackBrowserConfig 设置降级浏览器配置（当Unmarshal失败时使用）
func (l *Loader) setFallbackBrowserConfig(v *viper.Viper) {
	l.browserConfig.Basic.UserAgent = v.GetString("basic.user_agent")
	l.browserConfig.Basic.AcceptLanguage = v.GetString("basic.accept_language")
	l.browserConfig.Basic.Timezone = v.GetString("basic.timezone")
	l.browserConfig.Basic.Platform = v.GetString("basic.platform")
	l.browserConfig.Basic.Vendor = v.GetString("basic.vendor")
	l.browserConfig.Basic.Product = v.GetString("basic.product")

	l.browserConfig.Screen.Width = v.GetInt("screen.width")
	l.browserConfig.Screen.Height = v.GetInt("screen.height")
	l.browserConfig.Screen.AvailWidth = v.GetInt("screen.avail_width")
	l.browserConfig.Screen.AvailHeight = v.GetInt("screen.avail_height")
	l.browserConfig.Screen.ColorDepth = v.GetInt("screen.color_depth")
	l.browserConfig.Screen.PixelDepth = v.GetInt("screen.pixel_depth")
	l.browserConfig.Screen.DevicePixelRatio = v.GetFloat64("screen.device_pixel_ratio")

	l.browserConfig.Hardware.CPUCores = v.GetInt("hardware.cpu_cores")
	l.browserConfig.Hardware.DeviceMemory = v.GetInt("hardware.device_memory")
	l.browserConfig.Hardware.MaxTouchPoints = v.GetInt("hardware.max_touch_points")
	l.browserConfig.Hardware.VendorSub = v.GetString("hardware.vendor_sub")
	l.browserConfig.Hardware.ProductSub = v.GetString("hardware.product_sub")

	l.browserConfig.Canvas.EnableNoise = v.GetBool("canvas.enable_noise")
	l.browserConfig.Canvas.NoiseLevel = v.GetFloat64("canvas.noise_level")
	l.browserConfig.Canvas.BlockToDataURL = v.GetBool("canvas.block_toDataURL")

	l.browserConfig.WebGL.Vendor = v.GetString("webgl.vendor")
	l.browserConfig.WebGL.Renderer = v.GetString("webgl.renderer")
	l.browserConfig.WebGL.Version = v.GetString("webgl.version")
	l.browserConfig.WebGL.ShadingLanguageVersion = v.GetString("webgl.shading_language_version")

	l.browserConfig.Audio.EnableNoise = v.GetBool("audio.enable_noise")
	l.browserConfig.Audio.NoiseLevel = v.GetFloat64("audio.noise_level")

	l.browserConfig.WebRTC.BlockLocalIPLeak = v.GetBool("webrtc.block_local_ip_leak")
	l.browserConfig.WebRTC.FakePublicIP = v.GetString("webrtc.fake_public_ip")
	l.browserConfig.WebRTC.BlockDataChannels = v.GetBool("webrtc.block_data_channels")

	l.browserConfig.App.DefaultURL = v.GetString("app.default_url")
	l.browserConfig.App.WindowTitle = v.GetString("app.window_title")

	// 加载HTTP头部配置
	l.browserConfig.Headers.SecChUa = v.GetString("headers.sec_ch_ua")
	l.browserConfig.Headers.SecChUaMobile = v.GetString("headers.sec_ch_ua_mobile")
	l.browserConfig.Headers.SecChUaPlatform = v.GetString("headers.sec_ch_ua_platform")
	l.browserConfig.Headers.SecChUaFullVersionList = v.GetString("headers.sec_ch_ua_full_version_list")
	l.browserConfig.Headers.SecChUaArch = v.GetString("headers.sec_ch_ua_arch")
	l.browserConfig.Headers.SecChUaBitness = v.GetString("headers.sec_ch_ua_bitness")
	l.browserConfig.Headers.SecFetchDest = v.GetString("headers.sec_fetch_dest")
	l.browserConfig.Headers.SecFetchMode = v.GetString("headers.sec_fetch_mode")
	l.browserConfig.Headers.SecFetchSite = v.GetString("headers.sec_fetch_site")
	l.browserConfig.Headers.CacheControl = v.GetString("headers.cache_control")
	l.browserConfig.Headers.Pragma = v.GetString("headers.pragma")
	l.browserConfig.Headers.XSwCache = v.GetString("headers.x_sw_cache")

	l.browserConfig.Proxy.Mode = v.GetString("proxy.mode")
	l.browserConfig.Proxy.Url = v.GetString("proxy.url")
	l.browserConfig.Proxy.Username = v.GetString("proxy.username")
	l.browserConfig.Proxy.Password = v.GetString("proxy.password")
}

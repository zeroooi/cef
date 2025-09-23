// Package config 提供应用程序的配置管理功能
// 包含浏览器指纹伪装配置和访问控制配置的结构体定义
package config

// BrowserConfig 浏览器配置结构 - 完整的指纹伪装配置
type BrowserConfig struct {
	// 基础环境配置
	Basic struct {
		UserAgent      string `json:"user_agent"`
		AcceptLanguage string `json:"accept_language"`
		Timezone       string `json:"timezone"`
		Platform       string `json:"platform"`
		Vendor         string `json:"vendor"`
		Product        string `json:"product"`
	} `json:"basic"`

	// 屏幕和显示配置
	Screen struct {
		Width            int     `json:"width"`
		Height           int     `json:"height"`
		AvailWidth       int     `json:"avail_width"`
		AvailHeight      int     `json:"avail_height"`
		AvailTop         int     `json:"avail_top"`
		AvailLeft        int     `json:"avail_left"`
		ColorDepth       int     `json:"color_depth"`
		PixelDepth       int     `json:"pixel_depth"`
		DevicePixelRatio float64 `json:"device_pixel_ratio"`
	} `json:"screen"`

	// 硬件配置
	Hardware struct {
		CPUCores       int    `json:"cpu_cores"`
		DeviceMemory   int    `json:"device_memory"`
		MaxTouchPoints int    `json:"max_touch_points"`
		VendorSub      string `json:"vendor_sub"`
		ProductSub     string `json:"product_sub"`
	} `json:"hardware"`

	// Canvas指纹配置
	Canvas struct {
		EnableNoise    bool    `json:"enable_noise"`
		NoiseLevel     float64 `json:"noise_level"`
		BlockToDataURL bool    `json:"block_toDataURL"`
	} `json:"canvas"`

	// WebGL指纹配置
	WebGL struct {
		Vendor                 string   `json:"vendor"`
		Renderer               string   `json:"renderer"`
		Version                string   `json:"version"`
		ShadingLanguageVersion string   `json:"shading_language_version"`
		Extensions             []string `json:"extensions"`
	} `json:"webgl"`

	// 字体配置
	Fonts struct {
		AvailableFonts    []string `json:"available_fonts"`
		FontRandomization bool     `json:"font_randomization"`
	} `json:"fonts"`

	// 音频指纹配置
	Audio struct {
		EnableNoise bool    `json:"enable_noise"`
		NoiseLevel  float64 `json:"noise_level"`
	} `json:"audio"`

	// WebRTC配置
	WebRTC struct {
		BlockLocalIPLeak  bool   `json:"block_local_ip_leak"`
		FakePublicIP      string `json:"fake_public_ip"`
		BlockDataChannels bool   `json:"block_data_channels"`
	} `json:"webrtc"`

	// 插件配置
	Plugins struct {
		Enabled []struct {
			Name        string `json:"name"`
			Filename    string `json:"filename"`
			Description string `json:"description"`
		} `json:"enabled"`
	} `json:"plugins"`

	// 媒体设备配置
	MediaDevices struct {
		EnumerateDevicesNoise bool `json:"enumerate_devices_noise"`
		FakeDevices           []struct {
			Kind     string `json:"kind"`
			Label    string `json:"label"`
			DeviceId string `json:"deviceId"`
		} `json:"fake_devices"`
	} `json:"media_devices"`

	// 权限配置
	Permissions struct {
		Notifications string `json:"notifications"`
		Geolocation   string `json:"geolocation"`
		Camera        string `json:"camera"`
		Microphone    string `json:"microphone"`
	} `json:"permissions"`

	// 应用配置
	App struct {
		DefaultURL  string `json:"default_url"`
		WindowTitle string `json:"window_title"`
	} `json:"app"`

	// HTTP头部配置
	Headers struct {
		SecChUa                string `json:"sec_ch_ua"`
		SecChUaMobile          string `json:"sec_ch_ua_mobile"`
		SecChUaPlatform        string `json:"sec_ch_ua_platform"`
		SecChUaFullVersionList string `json:"sec_ch_ua_full_version_list"`
		SecChUaArch            string `json:"sec_ch_ua_arch"`
		SecChUaBitness         string `json:"sec_ch_ua_bitness"`
		SecFetchDest           string `json:"sec_fetch_dest"`
		SecFetchMode           string `json:"sec_fetch_mode"`
		SecFetchSite           string `json:"sec_fetch_site"`
		CacheControl           string `json:"cache_control"`
		Pragma                 string `json:"pragma"`
		XSwCache               string `json:"x_sw_cache"`
	} `json:"headers"`

	Proxy struct {
		Mode     string `json:"mode,omitempty"`
		Url      string `json:"url,omitempty"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
		Debug    bool   `json:"debug,omitempty"`
	} `json:"proxy"`
}

// WhitelistConfig 网站白名单配置结构
type WhitelistConfig struct {
	AllowedDomains    []string `json:"allowed_domains"`     // 允许访问的域名列表
	BlockedMessage    string   `json:"blocked_message"`     // 访问被阻止时的提示消息
	RedirectBlockedTo string   `json:"redirect_blocked_to"` // 被阻止时重定向的URL
}

// AppConfig 应用程序全局配置
type AppConfig struct {
	Browser   BrowserConfig   `json:"browser"`   // 浏览器指纹配置
	Whitelist WhitelistConfig `json:"whitelist"` // 白名单配置
}

type ExternalConfig struct {
	AegisAddr struct {
		Mode string `json:"mode"`
		Dev  string `json:"dev"`
		Long string `json:"long"`
		Pro  string `json:"pro"`
	} `json:"aegisAddr"`
}

type GetBrowserConfigResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ConfigMap map[string]*BrowserConfig `json:"config_map"` // ConfigMap
	} `json:"data"`
}

type GetWhitelistConfigResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ConfigMap map[string]*WhitelistConfig `json:"config_map"` // ConfigMap
	} `json:"data"`
}

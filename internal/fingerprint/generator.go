// Package fingerprint 指纹脚本生成器
// 根据配置参数动态生成JavaScript指纹伪装脚本
package fingerprint

import (
	"cef/internal/config"
	"fmt"
	"strings"
)

// Generator 指纹脚本生成器
type Generator struct {
	browserConfig *config.BrowserConfig
}

// NewGenerator 创建新的脚本生成器实例
func NewGenerator(browserConfig *config.BrowserConfig) *Generator {
	return &Generator{
		browserConfig: browserConfig,
	}
}

// GenerateBasicScript 根据配置文件参数创建完整的浏览器指纹伪装脚本
func (g *Generator) GenerateBasicScript() string {
	// 从配置中提取主语言
	primaryLanguage := g.extractPrimaryLanguage()
	// 从配置中生成语言数组
	languagesArray := g.generateLanguagesArray()

	return `

// 立即设置初始状态，避免时序问题
window.fingerprintOverridden = false;
window.fingerprintSuccess = {};
window.fingerprintData = {};

// 目标配置 - 精确匹配测试页面期望值
if (!window.__fingerprintConfig) {
    window.__fingerprintConfig = {
        userAgent: '` + g.browserConfig.Basic.UserAgent + `',
        platform: '` + g.browserConfig.Basic.Platform + `',
        hardwareConcurrency: ` + fmt.Sprintf("%d", g.browserConfig.Hardware.CPUCores) + `,
        language: '` + primaryLanguage + `',
        languages: ` + languagesArray + `,
        screenWidth: ` + fmt.Sprintf("%d", g.browserConfig.Screen.Width) + `,
        screenHeight: ` + fmt.Sprintf("%d", g.browserConfig.Screen.Height) + `,
        devicePixelRatio: ` + fmt.Sprintf("%.1f", g.browserConfig.Screen.DevicePixelRatio) + `
    };
    console.log('🎯 系统性修复配置:', window.__fingerprintConfig);
}
`
}

// GenerateAdvancedScript 创建高级指纹伪装脚本（Canvas、WebGL、音频等）
func (g *Generator) GenerateAdvancedScript() string {
	return `
(function() {
    // ========== Canvas指纹伪装 ==========
    if (` + fmt.Sprintf("%v", g.browserConfig.Canvas.EnableNoise) + `) {
        try {
            // Canvas 2D指纹伪装
            const originalGetImageData = CanvasRenderingContext2D.prototype.getImageData;
            CanvasRenderingContext2D.prototype.getImageData = function(sx, sy, sw, sh) {
                const imageData = originalGetImageData.apply(this, arguments);
                
                // 添加微小噪声
                const data = imageData.data;
                const noiseLevel = ` + fmt.Sprintf("%.6f", g.browserConfig.Canvas.NoiseLevel) + `;
                
                for (let i = 0; i < data.length; i += 4) {
                    const noise = (Math.random() - 0.5) * noiseLevel * 255;
                    data[i] = Math.max(0, Math.min(255, data[i] + noise));     // R
                    data[i + 1] = Math.max(0, Math.min(255, data[i + 1] + noise)); // G
                    data[i + 2] = Math.max(0, Math.min(255, data[i + 2] + noise)); // B
                }
                
                return imageData;
            };
            
            // Canvas toDataURL伪装
            const originalToDataURL = HTMLCanvasElement.prototype.toDataURL;
            HTMLCanvasElement.prototype.toDataURL = function() {
                const ctx = this.getContext('2d');
                if (ctx) {
                    // 添加不可见的像素噪声
                    const imageData = ctx.getImageData(0, 0, 1, 1);
                    const noise = Math.random() * 0.1;
                    imageData.data[0] += noise;
                    ctx.putImageData(imageData, 0, 0);
                }
                return originalToDataURL.apply(this, arguments);
            };
            
            console.log('Canvas指纹伪装完成');
            
        } catch(e) {
            console.error('Canvas伪装失败:', e);
        }
    }
    
    // ========== WebGL指纹伪装 ==========
    try { 
        // WebGL常量（直接使用数字值）
        const VENDOR = 0x1F00;
        const RENDERER = 0x1F01;
        const VERSION = 0x1F02;
        const SHADING_LANGUAGE_VERSION = 0x8B8C;
        
        // 目标WebGL配置 - 更全面的参数伪装
        const webglConfig = {
            [VENDOR]: '` + g.browserConfig.WebGL.Vendor + `',
            [RENDERER]: '` + g.browserConfig.WebGL.Renderer + `',
            [VERSION]: '` + g.browserConfig.WebGL.Version + `',
            [SHADING_LANGUAGE_VERSION]: '` + g.browserConfig.WebGL.ShadingLanguageVersion + `',
            // 额外的常见参数
            0x8B8A: 1, // MAX_VERTEX_ATTRIBS
            0x8DFB: 16, // MAX_TEXTURE_IMAGE_UNITS
            0x84E8: 16, // MAX_COMBINED_TEXTURE_IMAGE_UNITS
            0x8872: 4096, // MAX_TEXTURE_SIZE
            0x851C: 1024, // MAX_CUBE_MAP_TEXTURE_SIZE
            0x8073: 4, // SUBPIXEL_BITS
            0x80E9: 8, // SAMPLE_BUFFERS
            0x80EA: 4  // SAMPLES
        };
         
        // 保存原始方法
        const originalGetParameter = WebGLRenderingContext.prototype.getParameter;
        
        // 覆盖WebGL getParameter方法
        function overrideGetParameter(context) {
            context.getParameter = function(parameter) {
                console.log('🔍 WebGL getParameter调用，参数:', parameter);
                
                // 检查是否是我们要伪装的参数
                if (webglConfig.hasOwnProperty(parameter)) {
                    const fakeValue = webglConfig[parameter];
                    return fakeValue;
                }
                
                // 其他参数使用原始方法
                try {
                    const result = originalGetParameter.call(this, parameter);
                    console.log('🔍 原始参数', parameter, '返回:', result);
                    return result;
                } catch(e) {
                    return null;
                }
            };
        }
        
        // 覆盖WebGL 1.0
        if (typeof WebGLRenderingContext !== 'undefined') {
            WebGLRenderingContext.prototype.getParameter = function(parameter) {
                if (webglConfig.hasOwnProperty(parameter)) {
                    const fakeValue = webglConfig[parameter];
                    console.log('🎯 WebGL1 返回伪装值 [' + parameter + ']:', fakeValue);
                    return fakeValue;
                }
                return originalGetParameter.call(this, parameter);
            };
        }
        
        // 覆盖WebGL 2.0
        if (typeof WebGL2RenderingContext !== 'undefined') {
            WebGL2RenderingContext.prototype.getParameter = WebGLRenderingContext.prototype.getParameter;
        }
        
        // 立即测试
        setTimeout(function() {
            try {
                const testCanvas = document.createElement('canvas');
                const gl = testCanvas.getContext('webgl') || testCanvas.getContext('experimental-webgl');
                if (gl) {
                    console.log('🧪 WebGL伪装测试结果:');
                    console.log('  VENDOR (0x1F00):', gl.getParameter(0x1F00));
                    console.log('  RENDERER (0x1F01):', gl.getParameter(0x1F01));
                    console.log('  VERSION (0x1F02):', gl.getParameter(0x1F02));
                    console.log('  SHADING_LANGUAGE_VERSION (0x8B8C):', gl.getParameter(0x8B8C));
                } else {
                }
            } catch(testError) {
                console.error('WebGL测试失败:', testError);
            }
        }, 10);
        
        console.log('✅ WebGL指纹伪装完成');
        
    } catch(e) {
        console.error(' WebGL伪装失败:', e);
        console.error('错误堆栈:', e.stack);
    }
    
    // ========== 音频指纹伪装 ==========
    if (` + fmt.Sprintf("%v", g.browserConfig.Audio.EnableNoise) + `) {
        try {
            const AudioContext = window.AudioContext || window.webkitAudioContext;
            if (AudioContext) {
                const originalGetFloatFrequencyData = AnalyserNode.prototype.getFloatFrequencyData;
                AnalyserNode.prototype.getFloatFrequencyData = function(array) {
                    const result = originalGetFloatFrequencyData.apply(this, arguments);
                    
                    // 添加微小噪声
                    const noiseLevel = ` + fmt.Sprintf("%.6f", g.browserConfig.Audio.NoiseLevel) + `;
                    for (let i = 0; i < array.length; i++) {
                        array[i] += (Math.random() - 0.5) * noiseLevel;
                    }
                    
                    return result;
                };
            }
            
            // 音频格式伪装 - 标准化音频编解码器支持
            const originalCanPlayType = HTMLMediaElement.prototype.canPlayType;
            HTMLMediaElement.prototype.canPlayType = function(type) {
                // 标准化常见音频格式支持，避免独特性
                const commonFormats = {
                    'audio/ogg; codecs="vorbis"': 'probably',
                    'audio/ogg; codecs="opus"': 'probably', 
                    'audio/wav; codecs="1"': 'probably',
                    'audio/webm; codecs="vorbis"': 'probably',
                    'audio/webm; codecs="opus"': 'probably',
                    'audio/mp4': 'maybe',
                    'audio/mpeg': 'maybe',
                    'audio/flac': '',
                    'audio/wav': 'probably'
                };
                
                if (commonFormats.hasOwnProperty(type)) {
                    return commonFormats[type];
                }
                
                return originalCanPlayType.call(this, type);
            };
            
            console.log('音频指纹伪装完成');
            
        } catch(e) {
            console.error('音频伪装失败:', e);
        }
    }
    
    // ========== WebRTC IP泄露防护 ==========
    if (` + fmt.Sprintf("%v", g.browserConfig.WebRTC.BlockLocalIPLeak) + `) {
        try {
            const originalRTCPeerConnection = RTCPeerConnection;
            window.RTCPeerConnection = function(config) {
                if (config && config.iceServers) {
                    // 过滤STUN服务器防止IP泄露
                    config.iceServers = config.iceServers.filter(server => {
                        return !server.urls || !server.urls.toString().includes('stun');
                    });
                }
                return new originalRTCPeerConnection(config);
            };
            
            console.log('WebRTC IP泄露防护完成');
            
        } catch(e) {
            console.error('WebRTC防护失败:', e);
        }
    }
    
    // ========== WebSocket伪装和增强 ==========
    try {
        // 保存原始WebSocket构造函数
        const OriginalWebSocket = window.WebSocket;
        
        // 创建增强WebSocket构造函数
        window.WebSocket = function(url, protocols) { 
            // 创建原始WebSocket实例
            const ws = new OriginalWebSocket(url, protocols);
            
            // 添加错误处理和重连机制
            const originalOnError = ws.onerror;
            ws.onerror = function(event) {
                
                // 调用原始错误处理函数
                if (originalOnError) {
                    originalOnError.call(this, event);
                }
            };
            
            // 添加连接成功日志
            const originalOnOpen = ws.onopen;
            ws.onopen = function(event) {
                if (originalOnOpen) {
                    originalOnOpen.call(this, event);
                }
            };
            
            return ws;
        };
        
        // 保持WebSocket的原型链
        window.WebSocket.prototype = OriginalWebSocket.prototype;
        window.WebSocket.CONNECTING = OriginalWebSocket.CONNECTING;
        window.WebSocket.OPEN = OriginalWebSocket.OPEN;
        window.WebSocket.CLOSING = OriginalWebSocket.CLOSING;
        window.WebSocket.CLOSED = OriginalWebSocket.CLOSED;
        
        console.log('WebSocket伪装和增强完成');
        
    } catch(e) {
        console.error('WebSocket伪装失败:', e);
    }
    
    // ========== 字体指纹伪装 ==========
    try {
        // 字体检测伪装
        const availableFonts = ` + fmt.Sprintf("%q", g.browserConfig.Fonts.AvailableFonts) + `;
        
        // 拦截字体测量方法
        const originalMeasureText = CanvasRenderingContext2D.prototype.measureText;
        CanvasRenderingContext2D.prototype.measureText = function(text) {
            const result = originalMeasureText.apply(this, arguments);
            
            if (` + fmt.Sprintf("%v", g.browserConfig.Fonts.FontRandomization) + `) {
                // 添加微小的随机变化
                result.width += (Math.random() - 0.5) * 0.1;
            }
            
            return result;
        };
        
        console.log('字体指纹伪装完成');
        
    } catch(e) {
        console.error('字体伪装失败:', e);
    }
    
    // ========== 邮箱输入校验功能 ==========
    try {
        // 预设的邮箱白名单
        const allowedEmails = ["wuyan@yt-hsuanyuen.com", "suyunfei@hsuanyuen.com"];
		const pwdInput = document.getElementsByClassName('ace-input');
		pwdInput[1].disabled = true;
        
        // 显示提示信息的函数
        function showEmailWarning(message) {
            // 创建或更新提示元素
            let warningElement = document.getElementById('email-warning');
            if (!warningElement) {
                warningElement = document.createElement('div');
                warningElement.id = 'email-warning';
                warningElement.style.position = 'fixed';
                warningElement.style.top = '20px';
                warningElement.style.left = '50%';
                warningElement.style.transform = 'translateX(-50%)';
                warningElement.style.background = '#ff5555';
                warningElement.style.color = 'white';
                warningElement.style.padding = '10px 20px';
                warningElement.style.borderRadius = '4px';
                warningElement.style.zIndex = '10000';
                warningElement.style.fontSize = '14px';
                warningElement.style.boxShadow = '0 2px 10px rgba(0,0,0,0.2)';
                document.body.appendChild(warningElement);
            }
            
            warningElement.textContent = message;
            warningElement.style.display = 'block';
            
            // 3秒后自动隐藏
            setTimeout(() => {
                warningElement.style.display = 'none';
            }, 3000);
        }
        
        // 校验邮箱是否在白名单内
        function isEmailAllowed(email) {
            return allowedEmails.includes(email);
        }
        
        // 清空邮箱输入框
        function clearEmailInput(input) {
    		pwdInput[1].disabled = true;
			const button = document.getElementsByClassName('ace-ui-btn');
    		button[0].disabled = true;
        }
        
        // 校验并处理邮箱输入
        function validateAndHandleEmail(input) {
            if (!input) return false;
            
            const email = input.value.trim();
            if (email && !isEmailAllowed(email)) {
                clearEmailInput(input);
                showEmailWarning('邮箱不在允许列表中，请使用预设邮箱');
                return false;
            } else if (email && isEmailAllowed(email)) {
				pwdInput[1].disabled = false;
			}
            return true;
        }
        
        // 页面加载完成后设置事件监听器
        function setupEmailValidation() {
            // 查找邮箱输入框
            const emailInputs = document.querySelectorAll('input[name="email"][type="text"]');
            emailInputs.forEach(function(emailInput) {
                // 失焦事件处理
                emailInput.addEventListener('blur', function() {
                    validateAndHandleEmail(this);
                });
            });
            
            // 查找所有表单并拦截提交事件
            const forms = document.querySelectorAll('form');
            forms.forEach(function(form) {
                form.addEventListener('submit', function(e) {
                    // 查找表单中的邮箱输入框
                    const emailInput = form.querySelector('input[name="email"][type="text"]');
                    if (emailInput && !validateAndHandleEmail(emailInput)) {
                        e.preventDefault(); // 阻止表单提交
                        return false;
                    }
                });
            });
            
            // 也直接查找登录按钮并绑定事件
            const loginButtons = document.querySelectorAll('button[type="submit"], input[type="submit"]');
            loginButtons.forEach(function(button) {
                button.addEventListener('click', function(e) {
                    // 查找页面中的邮箱输入框
                    const emailInput = document.querySelector('input[name="email"][type="text"]');
                    if (emailInput && !validateAndHandleEmail(emailInput)) {
                        e.preventDefault(); // 阻止表单提交
                        e.stopPropagation(); // 阻止事件冒泡
                        return false;
                    }
                });
            });
        }
        
        // 等待DOM加载完成
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', setupEmailValidation);
        } else {
            // DOM已经加载完成
            setupEmailValidation();
        }
        
        // 对于动态加载的页面，使用MutationObserver监听DOM变化
        const observer = new MutationObserver(function(mutations) {
            let shouldSetup = false;
            
            mutations.forEach(function(mutation) {
                // 检查是否有新增节点
                mutation.addedNodes.forEach(function(node) {
                    if (node.nodeType === 1) { // 元素节点
                        // 检查是否是邮箱输入框或表单元素
                        if ((node.matches && (node.matches('input[name="email"][type="text"]') || 
                                              node.matches('form') || 
                                              node.matches('button[type="submit"], input[type="submit"]'))) ||
                            (node.querySelectorAll && 
                             (node.querySelectorAll('input[name="email"][type="text"]').length > 0 ||
                              node.querySelectorAll('form').length > 0 ||
                              node.querySelectorAll('button[type="submit"], input[type="submit"]').length > 0))) {
                            shouldSetup = true;
                        }
                    }
                });
            });
            
            if (shouldSetup) {
                // 延迟执行以确保DOM完全加载
                setTimeout(setupEmailValidation, 100);
            }
        });
        
        // 开始观察DOM变化
        observer.observe(document.body, {
            childList: true,
            subtree: true
        });
        
        console.log('邮箱输入校验功能已启用');
    } catch(e) {
        console.error('邮箱输入校验功能初始化失败:', e);
    }
    
    console.log('高级指纹伪装完成');
    
})();
`
}

// UpdateConfig 更新配置（运行时热更新）
func (g *Generator) UpdateConfig(newConfig *config.BrowserConfig) {
	g.browserConfig = newConfig
}

// GetConfigSummary 获取当前配置摘要（用于调试）
func (g *Generator) GetConfigSummary() map[string]interface{} {
	return map[string]interface{}{
		"user_agent":     g.browserConfig.Basic.UserAgent,
		"timezone":       g.browserConfig.Basic.Timezone,
		"screen_size":    fmt.Sprintf("%dx%d", g.browserConfig.Screen.Width, g.browserConfig.Screen.Height),
		"canvas_noise":   g.browserConfig.Canvas.EnableNoise,
		"audio_noise":    g.browserConfig.Audio.EnableNoise,
		"webrtc_blocked": g.browserConfig.WebRTC.BlockLocalIPLeak,
		"cpu_cores":      g.browserConfig.Hardware.CPUCores,
		"device_memory":  g.browserConfig.Hardware.DeviceMemory,
	}
}

// extractPrimaryLanguage 从AcceptLanguage配置中提取主语言
func (g *Generator) extractPrimaryLanguage() string {
	acceptLang := g.browserConfig.Basic.AcceptLanguage
	if acceptLang == "" {
		return "zh-CN" // 默认值
	}

	// 提取第一个语言标签
	languages := strings.Split(acceptLang, ",")
	if len(languages) > 0 {
		primaryLang := strings.TrimSpace(languages[0])
		// 移除质量值，如 "zh-CN;q=0.9" -> "zh-CN"
		if strings.Contains(primaryLang, ";") {
			primaryLang = strings.Split(primaryLang, ";")[0]
		}
		return strings.TrimSpace(primaryLang)
	}

	return "zh-CN" // 默认值
}

// generateLanguagesArray 生成语言数组的JavaScript代码
func (g *Generator) generateLanguagesArray() string {
	languages := strings.Split(g.browserConfig.Basic.AcceptLanguage, ",")
	var jsArray []string

	for _, lang := range languages {
		// 清理语言标签（移除质量值，如 "zh;q=0.9" -> "zh"）
		lang = strings.TrimSpace(lang)
		if strings.Contains(lang, ";") {
			lang = strings.Split(lang, ";")[0]
		}
		jsArray = append(jsArray, "'"+lang+"'")
	}

	return "[" + strings.Join(jsArray, ", ") + "]"
}

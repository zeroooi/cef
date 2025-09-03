// HTTP头部伪装脚本 - 智能修改请求头，避免与CEF层重复
// 专注于JavaScript层可控的头部，避免与CEF层冲突

(function() {
    'use strict';
    
    console.log('🌐 启动智能HTTP头部伪装（防重复版）...');
    
    // 自动提取User-Agent中的平台信息以确保Sec-Ch-Ua-Platform与其一致
    const userAgentString = navigator.userAgent;
    let platformValue = '"Windows"'; // 默认值
    
    if (userAgentString.indexOf('Windows') !== -1) {
        platformValue = '"Windows"';
    } else if (userAgentString.indexOf('Macintosh') !== -1) {
        platformValue = '"macOS"';
    } else if (userAgentString.indexOf('Linux') !== -1) {
        platformValue = '"Linux"';
    } else if (userAgentString.indexOf('Android') !== -1) {
        platformValue = '"Android"';
    } else if (userAgentString.indexOf('iPhone') !== -1 || userAgentString.indexOf('iPad') !== -1) {
        platformValue = '"iOS"';
    }

    // 强制设置 navigator.platform 以确保与 User-Agent 一致
    const platformRaw = platformValue.replace(/"/g, '');
    let navigatorPlatform = 'Win32';
    
    if (platformRaw === 'Windows') {
        navigatorPlatform = 'Win32';
    } else if (platformRaw === 'macOS') {
        navigatorPlatform = 'MacIntel';
    } else if (platformRaw === 'Linux') {
        navigatorPlatform = 'Linux x86_64';
    } else if (platformRaw === 'Android') {
        navigatorPlatform = 'Android';
    } else if (platformRaw === 'iOS') {
        navigatorPlatform = 'iPhone';
    }
    
    try {
        Object.defineProperty(navigator, 'platform', {
            get: function() {
                return navigatorPlatform;
            }
        });
    } catch (e) {
        console.warn(' 无法覆盖 navigator.platform:', e);
    }
    
    // 智能头部配置 - 只设置JavaScript层可以有效控制的头部
    // 避免与CEF层重复设置，防止头部值重复问题
    const jsControlledHeaders = {
        // 这些头部由JavaScript层专门处理
        // 注意：accept-language 由CEF层和动态脚本处理，避免在此硬编码
        'cache-control': 'no-cache',
        'pragma': 'no-cache',
        'x-sw-cache': '7'
    };
    
    // CEF层处理的头部（JavaScript不要重复设置）
    const cefControlledHeaders = [
        'sec-ch-ua-platform',  // CEF层根据User-Agent动态设置
        'user-agent'           // CEF层设置，避免冲突
    ];

    // ===== 强制覆盖头部方法 =====
    
    // 1. 智能检查函数 - 避免重复设置CEF控制的头部
    function shouldSkipHeader(headerName) {
        const lowerName = headerName.toLowerCase();
        return cefControlledHeaders.some(cefHeader => cefHeader.toLowerCase() === lowerName);
    }
    
    // 2. 劫持XMLHttpRequest原型的setRequestHeader方法
    const originalSetRequestHeader = XMLHttpRequest.prototype.setRequestHeader;
    XMLHttpRequest.prototype.setRequestHeader = function(name, value) {
        const lowerName = name.toLowerCase();
        
        // 跳过CEF控制的头部，避免冲突
        if (shouldSkipHeader(name)) {
            return; // 不调用原始方法
        }
        
        // 其他头部正常设置
        return originalSetRequestHeader.apply(this, arguments);
    };
    
    // 2. 拦截XMLHttpRequest
    if (window.XMLHttpRequest) {
        const originalOpen = XMLHttpRequest.prototype.open;
        const originalSend = XMLHttpRequest.prototype.send;
        
        XMLHttpRequest.prototype.open = function(method, url, async, user, password) {
            this._method = method;
            this._url = url;
            this._requestHeaders = {};
            return originalOpen.apply(this, arguments);
        };
        
        XMLHttpRequest.prototype.send = function(data) {

            // 只添加JavaScript层控制的头部，避免与CEF重复
            Object.keys(jsControlledHeaders).forEach(headerName => {
                if (!shouldSkipHeader(headerName)) {
                    try {
                        // 强化重复检测：检查现有头部值是否已包含目标值
                        const targetValue = jsControlledHeaders[headerName];
                        let shouldSet = true;
                        
                        // 获取当前已存在的头部值（如果有的话）
                        if (this._requestHeaders && this._requestHeaders[headerName.toLowerCase()]) {
                            console.log(`🔍 ${headerName} 已在记录中，跳过设置`);
                            shouldSet = false;
                        }
                        
                        if (shouldSet) {
                            this.setRequestHeader(headerName, targetValue);

                            // 记录已设置的头部
                            if (!this._requestHeaders) this._requestHeaders = {};
                            this._requestHeaders[headerName.toLowerCase()] = targetValue;
                        }
                    } catch (e) {
                        console.warn(`无法设置头部 ${headerName}:`, e);
                    }
                } else {
                    console.log(` 跳过CEF控制的头部 ${headerName}`);
                }
            });
            
            return originalSend.apply(this, arguments);
        };
    }
    
    // 3. 拦截Fetch API
    if (window.fetch) {
        const originalFetch = window.fetch;
        
        window.fetch = function(url, options = {}) {
            // 确保options和headers对象存在
            options = options || {};
            options.headers = options.headers || {};
            
            // 智能处理headers，避免与CEF层冲突
            let headersObj = options.headers;
            if (options.headers instanceof Headers) {
                headersObj = {};
                for (let [key, value] of options.headers.entries()) {
                    // 跳过CEF控制的头部
                    if (!shouldSkipHeader(key)) {
                        headersObj[key] = value;
                    } else {
                        console.log(`🔒 Fetch跳过CEF控制的头部 ${key}`);
                    }
                }
            }
            
            // 只添加JavaScript层控制的头部，强化重复检测
            Object.keys(jsControlledHeaders).forEach(headerName => {
                if (!shouldSkipHeader(headerName)) {
                    // 强化检查：检查是否已存在该头部
                    const existingHeader = Object.keys(headersObj).find(key => 
                        key.toLowerCase() === headerName.toLowerCase()
                    );
                    
                    if (!existingHeader) {
                        const targetValue = jsControlledHeaders[headerName];
                        headersObj[headerName] = targetValue;
                        console.log(` Fetch添加JS控制的头部 ${headerName}: ${targetValue}`);
                    } else {
                        console.log(` Fetch跳过已存在的头部 ${headerName}: ${headersObj[existingHeader]}`);
                    }
                } else {
                    console.log(` Fetch跳过CEF控制的头部 ${headerName}`);
                }
            });
            
            // 更新options
            options.headers = headersObj;
            
            return originalFetch.call(this, url, options);
        };
        
        console.log(' Fetch API头部拦截已设置');
    }
    
    // 4. 模拟navigator.userAgentData (Chrome 90+特性)
    if (!navigator.userAgentData) {
        // 从平台值提取不带引号的平台名
        const platformName = platformValue.replace(/"/g, '');
        
        Object.defineProperty(navigator, 'userAgentData', {
            get: function() {
                return {
                    brands: [
                        { brand: "Not_A Brand", version: "8" },
                        { brand: "Chromium", version: "120" },
                        { brand: "Google Chrome", version: "120" }
                    ],
                    mobile: false,
                    platform: platformName,
                    getHighEntropyValues: function(hints) {
                        return Promise.resolve({
                            platform: platformName,
                            platformVersion: platformName === "Windows" ? "10.0.0" : "15.0.0",
                            architecture: "x86",
                            model: "",
                            uaFullVersion: "120.0.0.0",
                            brands: this.brands
                        });
                    }
                };
            },
            enumerable: true,
            configurable: true
        });
        
    } else {
        // 如果已存在，尝试覆盖其平台值
        try {
            const platformName = platformValue.replace(/"/g, '');
            const originalUserAgentData = navigator.userAgentData;
            
            Object.defineProperty(navigator, 'userAgentData', {
                get: function() {
                    const data = {
                        brands: originalUserAgentData.brands,
                        mobile: originalUserAgentData.mobile,
                        platform: platformName,
                        getHighEntropyValues: function(hints) {
                            return Promise.resolve({
                                platform: platformName,
                                platformVersion: platformName === "Windows" ? "10.0.0" : "15.0.0",
                                architecture: "x86",
                                model: "",
                                uaFullVersion: "120.0.0.0",
                                brands: originalUserAgentData.brands
                            });
                        }
                    };
                    return data;
                }
            });
            
            console.log(' 已覆盖现有的navigator.userAgentData，平台设置为:', platformName);
        } catch (e) {
            console.warn(' 无法覆盖现有的navigator.userAgentData:', e);
        }
    }
    
    // 5. 监控所有网络请求
    try {
        const observer = new PerformanceObserver((list) => {
            list.getEntries().forEach((entry) => {
                if (entry.entryType === 'navigation' || entry.entryType === 'resource') {
                    console.log(` 网络请求: ${entry.name}`);
                }
            });
        });
        
        observer.observe({ entryTypes: ['navigation', 'resource'] });
        console.log(' 网络请求监控已启用');
    } catch (e) {
        console.warn(' PerformanceObserver不支持:', e);
    }
    
    // 6. 定期检查并强制更新头部（额外的保障措施）
    setInterval(() => {
        // 检查并强制更新navigator.platform
        try {
            if (navigator.platform !== navigatorPlatform) {
                Object.defineProperty(navigator, 'platform', {
                    get: function() {
                        return navigatorPlatform;
                    }
                });
                console.log(' 定期检查已重新设置 navigator.platform 为:', navigatorPlatform);
            }
        } catch (e) {
            console.warn(' 定期检查无法更新 navigator.platform:', e);
        }
        
        // 检查并强制更新userAgentData
        try {
            const platformName = platformValue.replace(/"/g, '');
            if (navigator.userAgentData && navigator.userAgentData.platform !== platformName) {
                Object.defineProperty(navigator, 'userAgentData', {
                    get: function() {
                        return {
                            brands: [
                                { brand: "Not_A Brand", version: "8" },
                                { brand: "Chromium", version: "120" },
                                { brand: "Google Chrome", version: "120" }
                            ],
                            mobile: false,
                            platform: platformName,
                            getHighEntropyValues: function(hints) {
                                return Promise.resolve({
                                    platform: platformName,
                                    platformVersion: platformName === "Windows" ? "10.0.0" : "15.0.0",
                                    architecture: "x86",
                                    model: "",
                                    uaFullVersion: "120.0.0.0",
                                    brands: [
                                        { brand: "Not_A Brand", version: "8" },
                                        { brand: "Chromium", version: "120" },
                                        { brand: "Google Chrome", version: "120" }
                                    ]
                                });
                            }
                        };
                    }
                });
                console.log(' 定期检查已重新设置 userAgentData.platform 为:', platformName);
            }
        } catch (e) {
            console.warn(' 定期检查无法更新 userAgentData:', e);
        }
    }, 1000); // 每秒检查一次
    
    // 7. 使用MutationObserver监控DOM变化，确保头部设置
    const observer = new MutationObserver(function(mutations) {
        mutations.forEach(function(mutation) {
            // 检查是否有新的iframe被添加，需要重新应用伪装
            if (mutation.type === 'childList') {
                mutation.addedNodes.forEach(function(node) {
                    if (node.tagName === 'IFRAME') {
                        console.log(' 检测到新iframe，重新应用头部伪装');
                        // 在iframe加载完成后重新应用伪装
                        node.addEventListener('load', function() {
                            try {
                                const iframeDoc = node.contentDocument || node.contentWindow.document;
                                if (iframeDoc) {
                                    // 在iframe中重新定义navigator属性
                                    const iframeWindow = node.contentWindow;
                                    if (iframeWindow) {
                                        Object.defineProperty(iframeWindow.navigator, 'platform', {
                                            get: function() {
                                                return navigatorPlatform;
                                            }
                                        });
                                        
                                        // 重新定义userAgentData
                                        Object.defineProperty(iframeWindow.navigator, 'userAgentData', {
                                            get: function() {
                                                return {
                                                    brands: [
                                                        { brand: "Not_A Brand", version: "8" },
                                                        { brand: "Chromium", version: "120" },
                                                        { brand: "Google Chrome", version: "120" }
                                                    ],
                                                    mobile: false,
                                                    platform: platformRaw,
                                                    getHighEntropyValues: function(hints) {
                                                        return Promise.resolve({
                                                            platform: platformRaw,
                                                            platformVersion: platformRaw === "Windows" ? "10.0.0" : "15.0.0",
                                                            architecture: "x86",
                                                            model: "",
                                                            uaFullVersion: "120.0.0.0",
                                                            brands: [
                                                                { brand: "Not_A Brand", version: "8" },
                                                                { brand: "Chromium", version: "120" },
                                                                { brand: "Google Chrome", version: "120" }
                                                            ]
                                                        });
                                                    }
                                                };
                                            }
                                        });
                                    }
                                }
                            } catch (e) {
                                console.warn(' iframe头部伪装失败:', e);
                            }
                        });
                    }
                });
            }
        });
    });
    
    // 开始观察DOM变化
    observer.observe(document, { childList: true, subtree: true });
    
})();
// 完整的浏览器指纹伪装脚本
// 用于隐藏自动化工具特征，模拟真实浏览器环境

(function() {
    'use strict';
    
    console.log('加载完整浏览器指纹伪装脚本...');
    
    try {
        // ========== 基础WebDriver检测隐藏 ==========
        
        // 1. 隐藏webdriver属性
        Object.defineProperty(navigator, 'webdriver', {
            get: () => undefined,
            configurable: true
        });
        
        // 2. 删除webdriver相关属性
        delete navigator.__proto__.webdriver;
        delete navigator.webdriver;
        
        // 3. 隐藏Chrome自动化相关属性
        if (window.chrome) {
            delete window.chrome.runtime;
            delete window.chrome.loadTimes;
            delete window.chrome.csi;
        }
        
        // 4. 隐藏常见的自动化检测全局变量
        const automationVars = [
            '_phantom', '__phantomas', 'callPhantom', 'callSelenium', 
            '_selenium', '__webdriver_evaluate', '__webdriver_script_func', 
            '__webdriver_script_fn', '_Selenium_IDE_Recorder', '_selenium', 
            'calledSelenium', '$cdc_asdjflasutopfhvcZLmcfl_', 
            '$chrome_asyncScriptInfo', '__$webdriverAsyncExecutor',
            'webdriver', '__webdriverFunc', '__lastWatirAlert',
            '__lastWatirConfirm', '__lastWatirPrompt', '_WEBDRIVER_ELEM_CACHE'
        ];
        
        automationVars.forEach(varName => {
            if (window.hasOwnProperty(varName)) {
                delete window[varName];
            }
        });
        
        // 5. 伪装chrome.app属性
        if (!window.chrome) {
            window.chrome = {};
        }
        
        if (!window.chrome.app) {
            Object.defineProperty(window.chrome, 'app', {
                get: () => ({
                    isInstalled: false,
                    InstallState: {
                        DISABLED: 'disabled',
                        INSTALLED: 'installed',
                        NOT_INSTALLED: 'not_installed'
                    },
                    RunningState: {
                        CANNOT_RUN: 'cannot_run',
                        READY_TO_RUN: 'ready_to_run',
                        RUNNING: 'running'
                    }
                }),
                configurable: true
            });
        }
        
        // 6. 伪装chrome.runtime
        if (!window.chrome.runtime) {
            Object.defineProperty(window.chrome, 'runtime', {
                get: () => ({
                    onConnect: undefined,
                    onMessage: undefined,
                    sendMessage: undefined
                }),
                configurable: true
            });
        }
        
        console.log('WebDriver检测隐藏完成');
        
    } catch (e) {
        console.error('WebDriver隐藏失败:', e);
    }
    
    try {
        // ========== 插件和权限伪装 ==========
        
        // 跳过插件列表伪装，因为navigator.plugins属性不可重新定义
        // 伪装MIME类型
        try {
            Object.defineProperty(navigator, 'mimeTypes', {
                get: () => ({
                    length: 0,
                    item: function() { return null; },
                    namedItem: function() { return null; }
                }),
                configurable: true
            });
        } catch(mimeError) {
            console.warn('MIME类型伪装失败:', mimeError.message);
        }
        
        console.log('插件和权限伪装完成（跳过了不可配置的plugins属性）');
        
    } catch (e) {
        console.error('插件伪装失败:', e);
    }
    
    try {
        // ========== 高级反检测技术 ==========
        
        // 1. 伪装iframe检测
        const originalcreateElement = document.createElement;
        document.createElement = function(tag) {
            const element = originalcreateElement.apply(this, arguments);
            
            if (tag.toLowerCase() === 'iframe') {
                // 隐藏可能的自动化标识
                setTimeout(() => {
                    if (element.contentWindow) {
                        try {
                            delete element.contentWindow.navigator.webdriver;
                        } catch (e) {}
                    }
                }, 0);
            }
            
            return element;
        };
        
        // 2. 伪装Date对象的getTimezoneOffset（防止时区检测）
        const originalGetTimezoneOffset = Date.prototype.getTimezoneOffset;
        Date.prototype.getTimezoneOffset = function() {
            // 返回上海时区的偏移量（-480分钟，即UTC+8）
            return -480;
        };
        
        // 3. 伪装console.debug（某些反爬工具会检测）
        if (!window.console.debug) {
            window.console.debug = function() {};
        }
        
        // 4. 伪装Error.captureStackTrace（V8特有）
        if (!Error.captureStackTrace) {
            Error.captureStackTrace = function() {};
        }
        
        console.log('高级反检测完成');
        
    } catch (e) {
        console.error('高级反检测失败:', e);
    }
    
    try {
        // ========== 行为模拟 ==========
        
        // 1. 模拟鼠标移动事件
        let mouseX = Math.random() * window.innerWidth;
        let mouseY = Math.random() * window.innerHeight;
        
        const simulateMouseMove = () => {
            mouseX += (Math.random() - 0.5) * 10;
            mouseY += (Math.random() - 0.5) * 10;
            
            mouseX = Math.max(0, Math.min(window.innerWidth, mouseX));
            mouseY = Math.max(0, Math.min(window.innerHeight, mouseY));
            
            const event = new MouseEvent('mousemove', {
                clientX: mouseX,
                clientY: mouseY,
                bubbles: true
            });
            
            document.dispatchEvent(event);
        };
        
        // 随机间隔模拟鼠标移动
        setInterval(simulateMouseMove, Math.random() * 5000 + 2000);
        
        // 2. 模拟页面焦点变化
        const simulateFocus = () => {
            if (Math.random() < 0.1) { // 10%概率触发
                window.dispatchEvent(new Event(Math.random() < 0.5 ? 'focus' : 'blur'));
            }
        };
        
        setInterval(simulateFocus, Math.random() * 10000 + 5000);
        
        console.log('行为模拟启动完成');
        
    } catch (e) {
        console.error('行为模拟失败:', e);
    }
    
    console.log('浏览器指纹伪装脚本加载完成');
    
})();
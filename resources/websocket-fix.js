// WebSocket深度修复脚本 - 优雅处理连接失败
// 该脚本通过静默处理WebSocket错误来避免无限重试

(function() {
    'use strict';
    
    console.log('🔧 启动WebSocket优雅错误处理...');
    
    // 保存原始WebSocket构造函数
    const OriginalWebSocket = window.WebSocket;
    
    // 创建优雅处理WebSocket实现
    function GracefulWebSocket(url, protocols) {
        console.log('🚀 创建WebSocket连接:', url);
        
        let ws;
        try {
            ws = new OriginalWebSocket(url, protocols);
        } catch (e) {
            return createMockWebSocket(url);
        }
        
        // 错误处理 - 避免无限重试
        const originalOnError = ws.onerror;
        ws.onerror = function(event) {

            // 不进行重试，直接静默处理
            this.readyState = OriginalWebSocket.CLOSED;
            
            // 调用原始错误处理器（如果存在）
            if (originalOnError && typeof originalOnError === 'function') {
                try {
                    originalOnError.call(this, event);
                } catch (e) {
                    console.warn('原始错误处理器执行失败:', e);
                }
            }
            
            // 阻止事件继续传播
            if (event && event.stopPropagation) {
                event.stopPropagation();
            }
        };
        
        // 连接成功处理
        const originalOnOpen = ws.onopen;
        ws.onopen = function(event) {
            if (originalOnOpen && typeof originalOnOpen === 'function') {
                originalOnOpen.call(this, event);
            }
        };
        
        // 关闭处理
        const originalOnClose = ws.onclose;
        ws.onclose = function(event) {
            console.log('🔌 WebSocket连接关闭:', event.code, event.reason || '未知原因');
            if (originalOnClose && typeof originalOnClose === 'function') {
                originalOnClose.call(this, event);
            }
        };
        
        return ws;
    }
    
    // 创建模拟WebSocket对象（用于连接失败时的降级处理）
    function createMockWebSocket(url) {
        const mockWs = {
            url: url,
            readyState: OriginalWebSocket.CLOSED,
            CONNECTING: OriginalWebSocket.CONNECTING,
            OPEN: OriginalWebSocket.OPEN,
            CLOSING: OriginalWebSocket.CLOSING,
            CLOSED: OriginalWebSocket.CLOSED,
            
            // 模拟方法
            send: function(data) {
                console.warn(' 模拟WebSocket发送（连接未建立）:', data);
            },
            
            close: function(code, reason) {
                console.log('🔌 模拟WebSocket关闭');
                this.readyState = OriginalWebSocket.CLOSED;
                if (this.onclose) {
                    this.onclose({ code: code || 1000, reason: reason || '手动关闭' });
                }
            },
            
            // 事件处理器
            onopen: null,
            onmessage: null,
            onerror: null,
            onclose: null,
            
            // 事件监听方法
            addEventListener: function(type, listener) {
                this['on' + type] = listener;
            },
            
            removeEventListener: function(type, listener) {
                if (this['on' + type] === listener) {
                    this['on' + type] = null;
                }
            }
        };
        
        // 模拟连接失败
        setTimeout(() => {
            if (mockWs.onerror) {
                mockWs.onerror({ type: 'error', target: mockWs });
            }
        }, 100);
        
        return mockWs;
    }
    
    // 保持原始WebSocket的属性和方法
    GracefulWebSocket.prototype = OriginalWebSocket.prototype;
    GracefulWebSocket.CONNECTING = OriginalWebSocket.CONNECTING;
    GracefulWebSocket.OPEN = OriginalWebSocket.OPEN;
    GracefulWebSocket.CLOSING = OriginalWebSocket.CLOSING;
    GracefulWebSocket.CLOSED = OriginalWebSocket.CLOSED;
    
    // 替换全局WebSocket
    window.WebSocket = GracefulWebSocket;
    
    // 阻止WebSocket错误冒泡到全局
    window.addEventListener('error', function(event) {
        if (event.target && event.target.constructor === OriginalWebSocket) {
            console.log('🛡️ 拦截WebSocket全局错误');
            event.preventDefault();
            event.stopPropagation();
            return false;
        }
    }, true);
    

})();
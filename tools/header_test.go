package main

import (
	"fmt"
	"strings"
)

// 模拟ICefStringMultiMap接口
type MockStringMultiMap struct {
	data map[string][]string
}

func NewMockStringMultiMap() *MockStringMultiMap {
	return &MockStringMultiMap{
		data: make(map[string][]string),
	}
}

func (m *MockStringMultiMap) GetSize() uint32 {
	return uint32(len(m.data))
}

func (m *MockStringMultiMap) GetKey(index uint32) string {
	i := 0
	for key := range m.data {
		if uint32(i) == index {
			return key
		}
		i++
	}
	return ""
}

func (m *MockStringMultiMap) GetValue(index uint32) string {
	i := 0
	for _, values := range m.data {
		if uint32(i) == index {
			if len(values) > 0 {
				return values[0]
			}
		}
		i++
	}
	return ""
}

func (m *MockStringMultiMap) Append(key, value string) {
	m.data[key] = append(m.data[key], value)
}

func (m *MockStringMultiMap) Clear() {
	m.data = make(map[string][]string)
}

func (m *MockStringMultiMap) Set(key, value string) {
	m.data[key] = []string{value}
}

// 测试头部去重功能
func testHeaderDeduplication() {
	fmt.Println("🧪 测试HTTP头部去重功能...")

	// 创建模拟的头部映射
	headerMap := NewMockStringMultiMap()

	// 添加一些重复的头部
	headerMap.Append("Accept-Language", "zh-CN,zh;q=0.9")
	headerMap.Append("Accept-Language", "zh-CN,zh;q=0.9") // 重复值
	headerMap.Append("Cache-Control", "no-cache")
	headerMap.Append("Cache-Control", "no-cache") // 重复值
	headerMap.Append("User-Agent", "Mozilla/5.0")
	headerMap.Append("Content-Type", "application/json")

	fmt.Println("📝 原始头部:")
	for i := uint32(0); i < headerMap.GetSize(); i++ {
		key := headerMap.GetKey(i)
		value := headerMap.GetValue(i)
		fmt.Printf("  %s: %s\n", key, value)
	}

	// 模拟去重过程
	deduplicated := deduplicateHeaders(headerMap)

	fmt.Println("\n✅ 去重后的头部:")
	for i := uint32(0); i < deduplicated.GetSize(); i++ {
		key := deduplicated.GetKey(i)
		value := deduplicated.GetValue(i)
		fmt.Printf("  %s: %s\n", key, value)
	}
}

// 头部去重函数（简化版）
func deduplicateHeaders(header *MockStringMultiMap) *MockStringMultiMap {
	preservedData := make(map[string]string)

	// 遍历所有数据，只保留每个键的第一个值
	size := header.GetSize()
	for i := uint32(0); i < size; i++ {
		key := header.GetKey(i)
		value := header.GetValue(i)

		// 如果键不存在，则添加
		if _, exists := preservedData[key]; !exists {
			preservedData[key] = value
		}
	}

	// 创建新的头部映射
	result := NewMockStringMultiMap()
	for key, value := range preservedData {
		result.Append(key, value)
	}

	return result
}

func main() {
	fmt.Println("🚀 HTTP头部去重测试工具")
	fmt.Println(strings.Repeat("=", 50))

	testHeaderDeduplication()

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("✅ 测试完成！")
}

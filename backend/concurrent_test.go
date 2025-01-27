package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestConcurrentBookPurchase(t *testing.T) {
	// 重置图书数据
	resp, err := http.Post("http://localhost:8080/admin/reset-books", "application/json", nil)
	if err != nil || resp.StatusCode != 200 {
		t.Fatal("重置图书数据失败")
	}

	// 获取初始图书列表
	resp, err = http.Get("http://localhost:8080/books")
	if err != nil {
		t.Fatal("获取图书列表失败:", err)
	}
	defer resp.Body.Close()

	var books []Book
	if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
		t.Fatal("解析图书数据失败:", err)
	}
	if len(books) == 0 {
		t.Fatal("图书列表为空")
	}

	// 选择第一本书进行测试
	testBook := books[0]
	initialStock := testBook.Stock

	// 设置并发购买的参数
	concurrentUsers := 10 // 并发用户数
	purchasePerUser := 2  // 每个用户购买的数量
	expectedStock := initialStock - (concurrentUsers * purchasePerUser)

	fmt.Printf("开始并发测试：\n")
	fmt.Printf("初始库存: %d\n", initialStock)
	fmt.Printf("并发用户数: %d\n", concurrentUsers)
	fmt.Printf("每用户购买数: %d\n", purchasePerUser)
	fmt.Printf("预期剩余库存: %d\n", expectedStock)

	// 使用WaitGroup来等待所有goroutine完成
	var wg sync.WaitGroup
	// 记录成功的购买次数
	var successCount atomic.Int32

	// 创建订单的函数
	createOrder := func() {
		defer wg.Done()

		order := Order{
			BookID:   testBook.ID,
			Quantity: purchasePerUser,
		}

		orderJSON, err := json.Marshal(order)
		if err != nil {
			t.Error("创建订单JSON失败:", err)
			return
		}

		resp, err := http.Post(
			"http://localhost:8080/order",
			"application/json",
			bytes.NewBuffer(orderJSON),
		)

		if err != nil {
			t.Error("发送订单请求失败:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			successCount.Add(1)
			fmt.Printf("订单创建成功\n")
		} else {
			var response map[string]string
			json.NewDecoder(resp.Body).Decode(&response)
			fmt.Printf("订单创建失败: %s\n", response["error"])
		}
	}

	// 启动并发购买
	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go createOrder()
	}

	// 等待所有购买完成
	wg.Wait()
	// 等待一秒确保数据库更新完成
	time.Sleep(time.Second)

	// 检查最终库存
	resp, err = http.Get("http://localhost:8080/books")
	if err != nil {
		t.Fatal("获取最终图书列表失败:", err)
	}
	defer resp.Body.Close()

	var finalBooks []Book
	if err := json.NewDecoder(resp.Body).Decode(&finalBooks); err != nil {
		t.Fatal("解析最终图书数据失败:", err)
	}

	// 找到测试的图书
	var finalStock int
	for _, book := range finalBooks {
		if book.ID == testBook.ID {
			finalStock = book.Stock
			break
		}
	}

	fmt.Printf("成功购买次数: %d\n", successCount)
	fmt.Printf("最终库存: %d\n", finalStock)

	if finalStock != expectedStock {
		t.Errorf("库存不符合预期: 期望 %d, 实际 %d", expectedStock, finalStock)
	} else {
		fmt.Printf("测试通过: 库存符合预期\n")
	}
}

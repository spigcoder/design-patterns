// Demo：书上的"装饰者模式在真实世界长什么样"——用 Go 标准库演示。
//
// 书上用的是 Java I/O：
//
//	new LineNumberInputStream(new BufferedInputStream(new FileInputStream("a.txt")))
//
// Go 里同一个模式无处不在，只是没人管它叫 Decorator，通常叫 wrapper / middleware：
//
//	io.LimitReader(bufio.NewReader(gzip.NewReader(file)))     ← io.Reader 装饰器
//	logging(auth(http.HandlerFunc(hello)))                    ← http.Handler 中间件
//
//	go run ./03-decorator/demo/stdlib
package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func main() {
	fmt.Println("=== 1. io.Reader 装饰器：同一个接口进，同一个接口出 ===")

	var r io.Reader = strings.NewReader("decorator pattern in go") // 被装饰者
	fmt.Printf("  原始            : %T\n", r)

	r = capsReader{r} // 自写装饰器：把输出转大写
	fmt.Printf("  +capsReader     : %T\n", r)

	r = bufio.NewReader(r) // 标准库装饰器：加缓冲
	fmt.Printf("  +bufio.Reader   : %T\n", r)

	r = io.LimitReader(r, 9) // 标准库装饰器：限制可读字节数
	fmt.Printf("  +io.LimitReader : %T\n", r)

	data, err := io.ReadAll(r)
	if err != nil {
		fmt.Println("读取失败：", err)
		return
	}
	fmt.Printf("  读出来          : %q\n", string(data))

	fmt.Println("\n=== 2. http.Handler 中间件：中间件就是装饰器 ===")

	hello := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello from "+r.URL.Path)
	})

	// 里层是 auth，外层是 logging；请求从外往里，响应从里往外出。
	handler := logging(auth(hello))
	fmt.Println("  组装顺序：logging(auth(hello))")

	call := func(label string, token string) {
		req := httptest.NewRequest(http.MethodGet, "/order", nil)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		fmt.Printf("  %s: %d %s %s\n", label, rec.Code, rec.Body.String(), afterAuth(rec))
	}

	call("没带令牌", "")
	call("令牌错误", "bad-token")
	call("令牌正确", "secret")

	fmt.Println("\n注意：hello 是 Handler，logging/auth 返回的也是 Handler ——")
	fmt.Println("和调料包住饮料一样：同类型进、同类型出，所以能无限套下去。")
}

// capsReader 是一个自写装饰器：不改读法，只把读到的字节变大写。
type capsReader struct{ r io.Reader }

func (c capsReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	for i := 0; i < n; i++ {
		if p[i] >= 'a' && p[i] <= 'z' {
			p[i] -= 'a' - 'A'
		}
	}
	return n, err
}

// logging 装饰器：打印一行日志，再交给下一层。
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("       [logging] %s -> %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// auth 装饰器：先检查令牌，不合格就短路，合格才交给下一层。
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// afterAuth 只是让输出更易读的辅助函数，与模式无关。
func afterAuth(rec *httptest.ResponseRecorder) string {
	if rec.Code == http.StatusUnauthorized {
		return "（auth 短路，没有进到 hello）"
	}
	return "（穿过 logging → auth → hello）"
}

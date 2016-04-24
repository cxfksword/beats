package color

import (
	"fmt"
	"testing"
	"time"
)

func Test_Color(t *testing.T) {
	fmt.Println(Color("color.Black", Black))
	fmt.Println(Color("color.Red", Red))
	fmt.Println(Color("color.Green", Green))
	fmt.Println(Color("color.Yellow", Yellow))
	fmt.Println(Color("color.Blue", Blue))
	fmt.Println(Color("color.Magenta", Magenta))
	fmt.Println(Color("color.Cyan", Cyan))
	fmt.Println(Color("color.White", White))
	fmt.Println(Color("color.Default", Default))
	// fmt.Println(Color("color.Gray", Gray))
}

func Test_Color_Output(t *testing.T) {
	fmt.Printf("[HTTP] %s %-15s %-4s %s\n%38s%s %s %s %s\n",
		time.Now().Format("15:04:05"),
		"192.168.10.233",
		"GET",
		"http://google.com", "←", StatusCodeColor(200), "text/html", "200KB",
		ContentColor("ddddddddddd"))
	fmt.Printf("[HTTP] %s %-15s %-4s %s\n%38s%s %s %s %s\n",
		time.Now().Format("15:04:05"),
		"192.168.105.231",
		"GET",
		"http://google.com", "←", StatusCodeColor(301), "text/html", "200KB",
		ContentColor("ddddddddddd"))
	fmt.Printf("[HTTP] %s %-15s %-4s %s\n%38s%s %s %s %s\n",
		time.Now().Format("15:04:05"),
		"192.168.10.233",
		"POST",
		"http://google.com", "←", StatusCodeColor(500), "text/html", "200KB",
		ContentColor("ddddddddddd"))
	fmt.Printf("%8s %-9s %-21s %-5s %-6s %s %4s %s\n",
		"[HTTP]",
		time.Now().Format("15:04:05"),
		fmt.Sprintf("%s:%d", "192.168.105.233", 8080),
		fmt.Sprintf("%dms", 43),
		fmt.Sprintf("%dKB", 13300/1000),
		StatusCodeColor(int(200)),
		"GET",
		"http://baidu.com/ddd")
	fmt.Printf("%8s %-9s %-21s %-5s %-6s %s\n",
		"[MYSQL]",
		time.Now().Format("15:04:05"),
		fmt.Sprintf("%s:%d", "127.0.0.1", 3060),
		fmt.Sprintf("%dms", 43),
		fmt.Sprintf("%dKB", 1500/1000),
		"select * from t")
	fmt.Printf("%8s %-9s %-21s %-5s %-6s %s %s\n",
		"[REDIS]",
		time.Now().Format("15:04:05"),
		fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		fmt.Sprintf("%dms", 43),
		fmt.Sprintf("%dKB", 153300/1000),
		"setnx key 1",
		"→ 25B")
	fmt.Printf("%8s %-9s %-21s %-5s %-6s %s\n",
		"[Thrift]",
		time.Now().Format("15:04:05"),
		fmt.Sprintf("%s:%d", "192.168.10.233", 9080),
		fmt.Sprintf("%dms", 43),
		fmt.Sprintf("%dKB", 15300/1000),
		fmt.Sprintf("%s%s", "hellowork(\"dddd\")", "\"ddddd\""))
	fmt.Printf("%8s %-9s %-21s %-5s %-6s %s %s %s\n",
		"[MC]",
		time.Now().Format("15:04:05"),
		fmt.Sprintf("%s:%d", "192.168.10.233", 11211),
		fmt.Sprintf("%dms", 43),
		fmt.Sprintf("%dKB", 153300/1000),
		"set",
		"key 2 39",
		"→ 25B")
}

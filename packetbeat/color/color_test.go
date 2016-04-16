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
	fmt.Println(Color("color.Gray", Gray))
	fmt.Println(Color("color.LightRed", LightRed))
	fmt.Println(Color("color.LightGreen", LightGreen))
	fmt.Println(Color("color.LightYellow", LightYellow))
	fmt.Println(Color("color.LightBlue", LightBlue))
	fmt.Println(Color("color.LightMagenta", LightMagenta))
	fmt.Println(Color("color.LightCyan", LightCyan))
	fmt.Println(Color("color.LightWhite", LightWhite))
}

func Test_Color_Http(t *testing.T) {
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
	fmt.Printf("[HTTP] %s %-15s %-4s %s\n%38s%s %s %s %s\n",
		time.Now().Format("15:04:05"),
		"192.168.10.233",
		"POST",
		"http://google.com", "←", StatusCodeColor(404), "text/html", "200KB",
		ContentColor(`

<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">

  <title>regexp - The Go Programming Language</title>

<link type="text/css" rel="stylesheet" href="/lib/godoc/style.css">

<link rel="search" type="application/opensearchdescription+xml" title="godoc" href="/opensearch.xml" />
`))
}

func Test_Color_Memcache(t *testing.T) {
	fmt.Printf("[MC] %s %-19s %s %s ← %s\n",
		time.Now().Format("15:04:05"),
		"192.168.10.233:6379",
		"GET",
		"dddd",
		"result")
	fmt.Printf("[MC] %s %-19s %s %s ← %s\n",
		time.Now().Format("15:04:05"),
		"127.0.0.1:6379",
		"GET",
		"dddd",
		"result")
}

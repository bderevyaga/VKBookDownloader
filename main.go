package main

import (
	"./vk"
	"./proxy"
)

func main() {
	proxy.Set()
	vk.Parse()
}

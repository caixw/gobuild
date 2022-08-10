// SPDX-License-Identifier: MIT

package main

func main() {
	println("test")
	exit := make(chan struct{})
	<-exit
}

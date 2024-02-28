// SPDX-FileCopyrightText: 2015-2024 caixw
//
// SPDX-License-Identifier: MIT

package main

func main() {
	println("test")
	exit := make(chan struct{})
	<-exit
}

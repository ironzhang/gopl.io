package main

import "fmt"

const (
	_ = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
	PiB
	EiB
	ZiB
	YiB
)

const (
	KB = 1000
	MB = 1000 * KB
	GB = 1000 * MB
	TB = 1000 * GB
	PB = 1000 * TB
	EB = 1000 * PB
	ZB = 1000 * EB
	YB = 1000 * ZB
)

func main() {
	fmt.Printf("KiB=%d\n", KiB)
	fmt.Printf("MiB=%d\n", MiB)
	fmt.Printf("GiB=%d\n", GiB)
	fmt.Printf("TiB=%d\n", TiB)
	fmt.Printf("PiB=%d\n", PiB)
	fmt.Printf("EiB=%d\n", EiB)
	fmt.Printf("YiB/ZiB=%v\n", YiB/ZiB)

	fmt.Printf("KB=%d\n", KB)
	fmt.Printf("MB=%d\n", MB)
	fmt.Printf("GB=%d\n", GB)
	fmt.Printf("TB=%d\n", TB)
	fmt.Printf("PB=%d\n", PB)
	fmt.Printf("EB=%d\n", EB)
	fmt.Printf("YB/ZB=%v\n", YB/ZB)
}

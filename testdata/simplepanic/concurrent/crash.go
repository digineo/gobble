package concurrent

func ILoveCrashing() {
	go concurrently()

	select {}
}

func concurrently() {
	panic("i died in a goroutine")
}

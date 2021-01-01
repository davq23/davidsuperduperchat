package channels

// OK checks if an element from a channel of boolean is true
func OK(done <-chan bool) bool {
	select {
	case ok := <-done:

		if ok {
			return true
		}
	}
	return false
}

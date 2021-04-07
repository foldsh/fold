package runtime

import "testing"

func TestAddrsAreRandom(t *testing.T) {
	addrs := []string{}
	for i := 0; i < 10; i++ {
		addrs = append(addrs, newAddr())
	}
	for i, addr := range addrs {
		for j, other := range addrs {
			if i != j {
				if addr == other {
					t.Fatalf("exepcted addrs to be random but found two the same")
				}
			}
		}
	}
}

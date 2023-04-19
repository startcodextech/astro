package os

import (
	"fmt"
	"os/user"
)

func IsRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Failed to get current user:", err)
		return false
	}

	return currentUser.Uid == "0"
}

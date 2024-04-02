package system

import (
	"fmt"
	"time"
)

func Log(args ...any) {
	timestamp := fmt.Sprintf("[%s]", time.Now().Format(time.RFC3339))
	line := append([]any{timestamp}, args...)
	fmt.Println(line...)
}

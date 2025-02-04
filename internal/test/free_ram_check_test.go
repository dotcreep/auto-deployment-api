package test

import (
	"fmt"
	"testing"

	"github.com/dotcreep/go-automate-deploy/internal/service/system"
)

func TestFreeRamCheck(t *testing.T) {
	percent, err := system.CheckFreeRam()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Free RAM: %d%%\n", percent)
}

func TestFreeRamCheckSecond(t *testing.T) {
	percent, err := system.CheckFreeRamSecond()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Usage RAM: %d%%\n", percent)
}

func TestFreeRamCheckThird(t *testing.T) {
	percent, err := system.CheckFreeRamThird()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Usage RAM: %d%%\n", percent)
}

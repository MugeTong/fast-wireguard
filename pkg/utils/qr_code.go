package utils

import (
	"fmt"
	"os"
	"github.com/mdp/qrterminal/v3"
)

/*
PrintQRCode prints a QR code to the terminal for the given string content.
*/
func PrintQRCode(prompt string, content string) {
	fmt.Printf("\n%s", prompt)
	fmt.Println("------------------------------------------------")

	// 配置二维码输出格式
	config := qrterminal.Config{
		Level:     qrterminal.L,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 1,
	}

	qrterminal.GenerateWithConfig(content, config)

	fmt.Println("------------------------------------------------")
	fmt.Println("End of QR Code")
}

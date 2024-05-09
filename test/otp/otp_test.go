package otp

import (
	"testing"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"bufio"
	"bytes"
	"encoding/base32"
	"fmt"
	"image/png"
	"os"
	"time"
)

func display(key *otp.Key, data []byte) {
	fmt.Printf("Issuer:       %s\n", key.Issuer())
	fmt.Printf("Account Name: %s\n", key.AccountName())
	fmt.Printf("Secret:       %s\n", key.Secret())
	fmt.Println("Writing PNG to qr-code.png....")
	os.WriteFile("qr-code.png", data, 0644)
	fmt.Println("")
	fmt.Println("Please add your TOTP to your OTP Application now!")
	fmt.Println("")
}

func promptForPasscode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Passcode: ")
	text, _ := reader.ReadString('\n')
	return text
}

// Demo function, not used in main
// Generates Passcode using a UTF-8 (not base32) secret and custom parameters
func GeneratePassCode(utf8string string) string {
	secret := base32.StdEncoding.EncodeToString([]byte(utf8string))
	passcode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	})
	if err != nil {
		panic(err)
	}
	return passcode
}

func TestMain(t *testing.T) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Knife",
		AccountName: "service@linuxcrypt.cn",
		// Secret:      []byte("A6MCIFEPOHLA44"),
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(key.Secret())
	v, e := totp.GenerateCode(key.Secret(), time.Now())
	fmt.Println(v, ":", e)

	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		panic(err)
	}
	png.Encode(&buf, img)

	// display the QR code to the user.
	display(key, buf.Bytes())
	fmt.Println(key.String())

	// Now Validate that the user's successfully added the passcode.
	// fmt.Println("Validating TOTP...")
	// passcode := promptForPasscode()
	// valid := totp.Validate(passcode, key.Secret())
	// if valid {
	// 	println("Valid passcode!")
	// 	os.Exit(0)
	// } else {
	// 	println("Invalid passcode!")
	// 	os.Exit(1)
	// }
}

func TestOtp(t *testing.T) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Knife",
		AccountName: "knife@linuxcrypt.cn",
		Period:      30,
		Secret:      []byte("iOO0rVqTCusye8hsDvpV"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%12s: %s\n", "Secret URL", key.URL())
	fmt.Printf("%12s: %s\n", "Key URL", key.String())

	now := time.Now()
	passcode, err := totp.GenerateCode(key.Secret(), now)
	if err != nil {
		panic(err)
	}
	fmt.Println("current passcode:", passcode)

	fmt.Println("验证")
	valid := totp.Validate(passcode, key.Secret())
	if valid {
		fmt.Println("Valid passcode!")
	} else {
		fmt.Println("Invalid passcode!")
	}
}

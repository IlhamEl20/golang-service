package utils

import (
	"log"
	"time"

	"github.com/dchest/captcha"
)

func init() {
	// Set the expiration time for captchas (e.g., 10 minute)
	captcha.SetCustomStore(captcha.NewMemoryStore(captcha.CollectNum, 10*time.Minute))
	log.Println("Captcha initialized")

}

// VerifyCaptcha verifies the captcha solution for the given captcha ID
func VerifyCaptcha(captchaID, solution string) bool {
	verified := captcha.VerifyString(captchaID, solution)
	return verified
}

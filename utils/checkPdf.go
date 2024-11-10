package utils

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func ValidatePDF(fileBytes []byte) error {
	// Gunakan pdfcpu untuk memvalidasi struktur file PDF
	err := api.Validate(bytes.NewReader(fileBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to validate PDF structure: %v", err)
	}

	// Log the file bytes
	// fmt.Printf("File bytes: %x\n", fileBytes)

	// Pencarian pola berbahaya menggunakan regex
	err = checkForMaliciousCode(fileBytes)
	if err != nil {
		return err
	}

	return nil
}

func checkForMaliciousCode(fileBytes []byte) error {
	// Daftar regex yang akan digunakan untuk mendeteksi kode berbahaya
	maliciousPatterns := []string{
		`(?i)javascript`,                   // Cari kata "JavaScript" (case insensitive)
		`(?i)<script.*?>.*?</script>`,      // Tag <script> HTML
		`/S\s*/JavaScript`,                 // Launch action JavaScript di PDF
		`/URI\s*\(`,                        // URI actions
		`(?i)<\?php.*?\?>`,                 // PHP code
		`(?i)import\s+(os|sys|subprocess)`, // Python import modules for os/system
		`(?i)(exec|shell_exec|system|passthru)\s*\(([^)]+)\)`, // Shell commands in execution context
		// `(?i)(\b(curl|wget|nc|bash|sh|chmod|chown|sudo)\b\s+[^\s])`, // Command injection usage
		`(?i)(SELECT|INSERT|UPDATE|DELETE|DROP|ALTER)\s+`, // SQL Injection
		`(?i)document\.write`,                             // JavaScript document.write
		`(?i)eval\s*\(`,                                   // eval function
		`(?i)onload\s*=`,                                  // onload event
		`(?i)onerror\s*=`,                                 // onerror event
		`(?i)onfocus\s*=`,                                 // onfocus event
		`(?i)onblur\s*=`,                                  // onblur event
		`(?i)onchange\s*=`,                                // onchange event
		`(?i)onsubmit\s*=`,                                // onsubmit event
		`(?i)onclick\s*=`,                                 // onclick event
		`(?i)ondblclick\s*=`,                              // ondblclick event
		`(?i)onkeydown\s*=`,                               // onkeydown event
		`(?i)onkeypress\s*=`,                              // onkeypress event
		`(?i)onkeyup\s*=`,                                 // onkeyup event

	}
	// Periksa setiap pola pada konten file PDF
	for _, pattern := range maliciousPatterns {
		matched, err := regexp.Match(pattern, fileBytes)
		if err != nil {
			return fmt.Errorf("failed to apply regex: %v", err)
		}
		if matched {
			return fmt.Errorf("malicious code detected matching pattern: %s", pattern)
		}
	}

	// Jika tidak ada pola mencurigakan yang ditemukan
	return nil
}

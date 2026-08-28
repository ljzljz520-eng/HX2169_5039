package importer

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

func Checksum(content string) string {
	digest := sha256.Sum256([]byte(content))
	return hex.EncodeToString(digest[:])
}

func VerifyChecksum(content, expected string) bool {
	if strings.TrimSpace(expected) == "" {
		return false
	}
	return Checksum(content) == expected
}

func ValidateAttachment(filename, content, expected string) error {
	if strings.TrimSpace(filename) == "" {
		return errors.New("filename is required")
	}
	if content == "" {
		return errors.New("content is required")
	}
	if expected != "" && !VerifyChecksum(content, expected) {
		return errors.New("checksum mismatch")
	}
	return nil
}

func Extension(filename string) string {
	index := strings.LastIndex(filename, ".")
	if index < 0 || index == len(filename)-1 {
		return ""
	}
	return strings.ToLower(filename[index+1:])
}

func SupportedExtension(filename string) bool {
	switch Extension(filename) {
	case "txt", "pdf", "docx":
		return true
	default:
		return false
	}
}

func NormalizeFilename(filename string) string {
	return strings.TrimSpace(strings.ReplaceAll(filename, "\\", "/"))
}

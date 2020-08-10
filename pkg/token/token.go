package token

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrMissingHeader = errors.New("The length of the `Authorization` header is zero.")
	ErrBasicRequired = errors.New("The Authorization should be basic type.")
)

// ParseRequest gets the Authorization from the header and parse it.
func ParseRequest(c *gin.Context) error {
	header := c.Request.Header.Get("Authorization")

	if len(header) == 0 {
		return ErrMissingHeader
	}

	var t string
	// Parse the header to get the token part.
	_, err := fmt.Sscanf(header, "Basic %s", &t)
	if err != nil {
		return err
	}

	sDec, err := base64.StdEncoding.DecodeString(t)
	if err != nil {
		return err
	}

	str := string(sDec)
	i := strings.Index(string(sDec), ":")
	if i < 0 {
		return ErrMissingHeader
	}

	c.Set("Sid", str[:i])
	c.Set("Password", str[i+1:])

	return nil
}

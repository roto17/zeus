package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/roto17/zeus/lib/sharedkeys"
)

func Coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return "" // return an empty string if all values are empty
}

func GetHeaderVarToString(value any, exists bool) string {

	if !exists {
		return ""
	}

	return value.(string)
}

func StringToInt(value string) int {

	output, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("%v", err)
	}

	return output
}

func DeviceNameString(userAgent string) string {
	// You can use a simple check or a more sophisticated library to parse the User-Agent
	// For example, let's just return the User-Agent string for simplicity
	if userAgent == "" {
		return "unknown device"
	}
	return userAgent
}

func GetTheOriginalIPAddressFromForwarded(IPS string) string {
	var originalIP string
	if IPS != "" {
		// X-Forwarded-For can contain multiple IPs, split by comma
		ipList := strings.Split(IPS, ",")
		originalIP = strings.TrimSpace(ipList[0]) // Take the first IP
	}
	return originalIP
}

func ExtractUserDetails(c *gin.Context) (jwt.MapClaims, error) {
	// user is already a map[string]interface{}, no need to assert it to jwt.MapClaims

	usr, _ := c.Get(sharedkeys.UserKey)

	claims, _ := usr.(jwt.MapClaims)

	return claims, nil
}

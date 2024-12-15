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

// GetCompanyIDFromGinClaims retrieves the company ID from the user's claims stored in Gin context.
func GetParamIDFromGinClaims(c *gin.Context, param string) uint {
	// Retrieve user from the Gin context
	user, exists := c.Get(sharedkeys.UserKey)
	if !exists {
		fmt.Printf("User claims do not exist\n")
		return 0
	}

	// Assert user to jwt.MapClaims
	claims, ok := user.(jwt.MapClaims)
	if !ok {
		fmt.Printf("Failed to assert User to jwt.MapClaims\n")
		return 0
	}

	// Retrieve and convert company_id from claims
	paramIDStr, ok := claims[param].(string)
	if !ok {
		fmt.Printf("%s claim is not a string\n", param)
		return 0
	}

	// Convert company_id string to uint
	paramID, err := strconv.Atoi(paramIDStr)
	if err != nil {
		fmt.Printf("Failed to convert %s from string to integer: %s\n", err, param)
		return 0
	}

	return uint(paramID)
}

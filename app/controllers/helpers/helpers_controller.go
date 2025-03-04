package helpers

import (
	"ekycapp/config"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Generate a Mob OTP string of specified length.
//
// @param length - The length of the string to be generated
//
// @return A string of random characters of specified length
func OtpGen(length int) string {

	bytes := make([]byte, length)

	for i := 0; i < length; i++ {
		bytes[i] = "0123456789"[rand.Intn(10)]
	}
	return string(bytes)
}

// ApxIdStringGen generates a random string of the specified length.
//
// @param length - The length of the string to generate. Must be a positive integer.
// @param chars - The characters to use.
//
// @return A randomly generated string of the specified length and character set.
func IdStringGen(length int, chars string, prefix string) (string, error) {
	// if length is negative panic
	if length <= 0 {
		return "", fmt.Errorf("length must be a positive integer : func IdStringGen")

	}

	if chars == "" {
		chars = "123456789ABCDEFGHIJKLMNPQRSTUVWXYZ" // Default characters
	}

	source := []rune(chars)
	var result string

	for i := 0; i < length; i++ {
		index := rand.Intn(len(source))
		result += string(source[index])
	}

	return prefix + result, nil
}

// generating random string for application_id
func ApplicationIdGen(length int, prefix, chars string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be a positive integer func ApplicationIdGen")
	}

	if chars == "" {
		chars = "123456789ABCDEFGHJKLMNPQRSTUVWXYZ" // Default characters
	}

	source := []rune(chars)

	var result string
	for i := 0; i < length; i++ {
		index := rand.Intn(len(source)) // Generate random index within source
		result += string(source[index])
	}

	return prefix + result, nil
}

// it is a replacement for Arr::get() laravel function which use to fix the offset issue
func GetValue(m map[string]interface{}, key string) interface{} {

	value, ok := m[key]
	if !ok || value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case int:
		return v
	case []interface{}:
		return v
	case map[string]interface{}:
		return v
	default:
		return ""
	}
}

// GenerateToken generates a token for the user
func GenerateToken(mobileNo, userType string) (string, error) {
	var secretKey any
	if userType == "kyc" {
		secretKey = []byte(config.Denv("KYC_TOKEN_SECRET"))
	}
	if userType == "admin" {
		secretKey = []byte(config.Denv("ADMIN_TOKEN_SECRET"))
	}
	if userType == "closure" {
		secretKey = []byte(config.Denv("CLOSURE_TOKEN_SECRET"))
	}
	if userType == "rekyc" {
		secretKey = []byte(config.Denv("REKYC_TOKEN_SECRET"))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"mobile": mobileNo,
		"type":   userType,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(secretKey)
}

// formating date to d/m/Y
func DateFormat(dateString, layout string) (string, error) {
	// layout := "02-01-2006"
	parse, err := time.Parse(layout, dateString)
	if err != nil {
		return "", err
	}
	newLayout := "02/01/2006"
	return parse.Format(newLayout), nil
}

func MergeMaps(m1, m2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m1 {
		result[k] = v
	}
	for k, v := range m2 {
		result[k] = v
	}
	return result
}

func GenerateFilenames(mobile, extension string) string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filenameNow := fmt.Sprintf("%s_%s%s", mobile, timestamp, "."+extension)

	return filenameNow
}

func HasAllKeys(defaultMap map[string]bool, childMap map[string]interface{}) bool {
	for key := range childMap {
		if _, exists := defaultMap[key]; !exists {
			return false
		}
	}

	for key := range defaultMap {
		if _, exists := childMap[key]; !exists {
			fmt.Println(key)
			return false
		}
	}
	return true
}

func HasSomeKeys(defaultMap map[string]bool, childMap map[string]interface{}) bool {
	for key := range childMap {
		if _, exists := defaultMap[key]; !exists {
			return false
		}
	}
	return true
}

// Assign time.Now() to a variable first and then take its address.
func TimePointer(t time.Time) *time.Time {
	return &t
}

func AddressFormat(inputString string, maxChars int) []string {
	// Split the input string by spaces to get words
	words := strings.Fields(inputString)
	var substrings []string
	var currentSubstring string

	for _, word := range words {
		// Append the current word to the current substring
		newSubstring := currentSubstring + word

		// Check if the length of the new substring is within the limit
		if utf8.RuneCountInString(newSubstring) <= maxChars {
			currentSubstring = newSubstring
		} else {
			// If the new substring exceeds the limit, store the current substring and start a new one
			substrings = append(substrings, currentSubstring)
			currentSubstring = word
		}
	}

	// Append any remaining substring
	if len(currentSubstring) > 0 {
		substrings = append(substrings, currentSubstring)
	}

	return substrings
}

func search() *gorm.DB {

	return config.DB
}

func DecodeAndParse(encoded string) map[string]interface{} {
	if encoded == "" {
		return nil
	}

	// Decode the Base64 string
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil
	}

	// // Parse the decoded JSON into a map
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(string(decodedBytes)), &result); err != nil {
		return nil
	}

	return result
}

// function to check if a key exists in a list
func Contains(list []string, key string) bool {
	for _, item := range list {
		if item == key {
			return true
		}
	}
	return false
}

// Convert to JSON string
func JsonToString(data interface{}) string {
	// Marshal the data to a JSON byte slice
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		// Return empty string on error
		return ""
	}
	// Convert to string and return
	return string(jsonBytes)
}

// function that accepts a struct and custom key-value pairs, returning a final object (map[string]interface{}) that combines the struct fields with the provided custom key-value pairs.
func StructToMapWithCf(inputStruct interface{}, customFields map[string]interface{}) (map[string]interface{}, error) {
	// Step 1: Marshal the struct into JSON
	structJSON, err := json.Marshal(inputStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct: " + err.Error())
	}

	// Step 2: Unmarshal the JSON into a map
	resultMap := make(map[string]interface{})
	if err := json.Unmarshal(structJSON, &resultMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal struct to map: " + err.Error())
	}

	// Step 3: Add custom key-value pairs to the map
	for key, value := range customFields {
		resultMap[key] = value
	}

	// Step 4: Return the final map
	return resultMap, nil
}

func ModelConvertString(model interface{}) (string, error) {
	jsonData, err := json.Marshal(model)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

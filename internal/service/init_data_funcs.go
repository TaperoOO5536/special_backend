package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/TaperoOO5536/special_backend/internal/models"
)

var (
	ErrFailedToParseInitData =     errors.New("failed to parse inti data")
	ErrUserNotFoundInInitData =    errors.New("user not fount in inti data")
	ErrFailedToUnmarshalUserData = errors.New("failed to unmarshal user data")
	ErrHashNotFound =              errors.New("hash not found in initData")
	ErrInitDataIsEmpty =           errors.New("init data is empty")
	ErrFailedToUnmarshalId =       errors.New("failed to unmarshar id")
)

type RawUser struct {
	ID          json.RawMessage `json:"id"`
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	UserName    string          `json:"username"`
	PhoneNumber string          `json:"phone_number,omitempty"`
}

func ParseInitData(initData string) (*models.User, error) {
	var user *models.User
	if initData == "" {
		return nil, ErrInitDataIsEmpty
	}

	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, ErrFailedToParseInitData
	}

	userData := values.Get("user")
	if userData == "" {
		return nil, ErrUserNotFoundInInitData
	}

	user = &models.User{}
	var rawUser RawUser
	err = json.Unmarshal([]byte(userData), &rawUser)
	if err != nil {
		return nil, ErrFailedToUnmarshalUserData
	}

	var idNum int64
	if err := json.Unmarshal(rawUser.ID, &idNum); err == nil {
		user.ID = strconv.FormatInt(idNum, 10)
	} else {
		var idStr string
		if err := json.Unmarshal(rawUser.ID, &idStr); err != nil {
			return nil, ErrFailedToUnmarshalId
		}
		user.ID = idStr
	}

	user.Name = rawUser.FirstName
	user.Surname = rawUser.LastName
	user.Nickname = rawUser.UserName

	return user, nil
}

func GetSecretKey(token string) []byte {
	h := sha256.New()
	h.Write([]byte("WebAppData"))
	secretKey := h.Sum(nil)
	h = hmac.New(sha256.New, secretKey)
	h.Write([]byte(token))
	key := h.Sum(nil)
	return key
}

func VerifyInitData(initData string, token string) (bool, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return false, ErrFailedToParseInitData
	}

	receivedHash := values.Get("hash")
	if receivedHash == "" {
		return false, ErrHashNotFound
	}

	values.Del("hash")
	var dataCheckStrings []string
	for key := range values {
		dataCheckStrings = append(dataCheckStrings, fmt.Sprintf("%s=%s", key, values.Get(key)))
	}
	sort.Strings(dataCheckStrings)
	dataCheckString := strings.Join(dataCheckStrings, "\n")

	secretKey := GetSecretKey(token)
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataCheckString))
	computedHash := hex.EncodeToString(h.Sum(nil))

	return computedHash == receivedHash, nil
}
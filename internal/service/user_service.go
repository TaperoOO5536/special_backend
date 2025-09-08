package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"github.com/TaperoOO5536/special_backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidInitData = errors.New("invalid init data")
	ErrFailedToParseInitData = errors.New("failed to parse inti data")
	ErrUserNotFoundInInitData = errors.New("user not fount in inti data")
	ErrFailedToUnmarshalUserData = errors.New("failed to unmarshal user data")
	ErrHashNotFound = errors.New("hash not found in initData")
)

type UserService struct {
	userRepo repository.UserRepository
	token    string
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserInfo(ctx context.Context, initData string) (*models.User, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, ErrInvalidInitData
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return nil, ErrFailedToParseInitData
	}

	userInfo, err := s.userRepo.GetUserInfo(ctx, user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrItemNotFound
		}
		return nil, err
	}

	return userInfo, nil
}

func ParseInitData(initData string) (*models.User, error) {
    var user *models.User
    values, err := url.ParseQuery(initData)
    if err != nil {
        return nil, ErrFailedToParseInitData
    }

    userData := values.Get("user")
    if userData == "" {
        return nil, ErrUserNotFoundInInitData
    }

    err = json.Unmarshal([]byte(userData), &user)
    if err != nil {
        return user, ErrFailedToUnmarshalUserData
    }

    return user, nil
}

func GetSecretKey(botToken string) []byte {
    h := sha256.New()
    h.Write([]byte("WebAppData"))
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
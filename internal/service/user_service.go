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
	"strconv"
	"strings"

	"github.com/TaperoOO5536/special_backend/internal/models"
	"github.com/TaperoOO5536/special_backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound =              errors.New("user not found")
	ErrInvalidInitData =           errors.New("invalid init data")
	ErrFailedToParseInitData =     errors.New("failed to parse inti data")
	ErrUserNotFoundInInitData =    errors.New("user not fount in inti data")
	ErrFailedToUnmarshalUserData = errors.New("failed to unmarshal user data")
	ErrHashNotFound =              errors.New("hash not found in initData")
	ErrInitDataIsEmpty =           errors.New("init data is empty")
	ErrUserExists =                errors.New("user already exists")
	ErrFailedToUnmarshalId =       errors.New("failed to unmarshar id")
)

type RawUser struct {
	ID          json.RawMessage `json:"id"`
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	UserName    string          `json:"username"`
	PhoneNumber string          `json:"phone_number,omitempty"`
}

type UserService struct {
	userRepo repository.UserRepository
	token    string
}

func NewUserService(userRepo repository.UserRepository, token string) *UserService {
	return &UserService{
		userRepo: userRepo,
		token: token,
	}
}

func (s *UserService) CreateUser(ctx context.Context, initData string) (error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return ErrInvalidInitData
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return err
	}

	existingUser, _ := s.userRepo.GetUserInfo(ctx, user.ID)
	if existingUser != nil {
		return ErrUserExists
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *UserService) GetUserInfo(ctx context.Context, initData string) (*models.User, error) {
	valid, err := VerifyInitData(initData, s.token)
	if err != nil || !valid {
		return nil, ErrInvalidInitData
	}

	user, err := ParseInitData(initData)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.userRepo.GetUserInfo(ctx, user.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return userInfo, nil
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

	user.Name =        rawUser.FirstName
	user.Surname =     rawUser.LastName
	user.Nickname =    rawUser.UserName
	user.PhoneNumber = rawUser.PhoneNumber

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
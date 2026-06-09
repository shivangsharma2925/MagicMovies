package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPService struct {
	redis *redis.Client
}

func NewOTPService(redisClient *redis.Client) *OTPService {
	return &OTPService{
		redis: redisClient,
	}
}

func (o *OTPService) GenerateSecureOTP() (string, error) {
	// The range size is 900000 (from 0 to 899999)
	max := big.NewInt(900000)

	// Securely pick a random number in [0, 899999]
	randomNum, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Shift range to [100000, 999999] to guarantee a 6-digit OTP
	otp := randomNum.Int64() + 100000

	return strconv.FormatInt(otp, 10), nil
}

func (o *OTPService) SaveVerificationOTP(ctx context.Context, UserID string, otp string) error {

	key := fmt.Sprintf("verify:%s", UserID)

	err :=  o.redis.Set(
		ctx,
		key,
		otp,
		10*time.Minute,
	).Err()

	return err
}

func (o *OTPService) GetVerificationOTP(ctx context.Context, UserID string,) (string, error) {

	key := fmt.Sprintf("verify:%s", UserID)

	return o.redis.Get(ctx, key).Result()
}

func (o *OTPService) DeleteVerificationOTP(ctx context.Context, UserID string,) error {

	key := fmt.Sprintf("verify:%s", UserID)

	return o.redis.Del(ctx, key).Err()
}

func (o *OTPService) IncrementAttempts(ctx context.Context, UserID string,) (int64, error) {

	key := fmt.Sprintf("verify_attempts:%s", UserID)

	count, err := o.redis.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	
	if count == 1 {
		o.redis.Expire(ctx, key, 10*time.Minute)
	}
	
	return count, nil
}

func (o *OTPService) ClearAttempts(ctx context.Context, UserID string,) error {

	key := fmt.Sprintf("verify_attempts:%s", UserID)

	return o.redis.Del(ctx, key).Err()
}

func (o *OTPService) IsResendCooldownActive(ctx context.Context, userID string) (bool, error) {

	key := fmt.Sprintf("verify_cooldown:%s", userID)

	exists, err := o.redis.Exists(ctx, key).Result()

	return exists > 0, err
}

func (o *OTPService) SetResendCooldown(ctx context.Context, userID string) error {

	key := fmt.Sprintf("verify_cooldown:%s", userID)

	return o.redis.Set(
		ctx,
		key,
		"1",
		time.Minute,
	).Err()
}

func (o *OTPService) SaveResetPasswordToken(ctx context.Context, userId, key string) error {
	redisKey := fmt.Sprintf("resetToken:%s", key)

	return o.redis.Set(ctx, redisKey, userId, 15*time.Minute).Err()
}

func (o *OTPService) GetResetPasswordToken(ctx context.Context, key string) (string, error) {
	redisKey := fmt.Sprintf("resetToken:%s", key)

	return o.redis.Get(ctx, redisKey).Result()
}

func (o *OTPService) DeleteResetPasswordToken(ctx context.Context, key string) error {
	redisKey := fmt.Sprintf("resetToken:%s", key)

	return o.redis.Del(ctx, redisKey).Err()
}

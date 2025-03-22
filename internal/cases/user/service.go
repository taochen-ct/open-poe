package user

import (
	"context"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sony/sonyflake"
	"go.uber.org/zap"
	"log"
	"open-poe/config"
	"open-poe/internal/compo"
	"open-poe/internal/pkg/request/user"
	"open-poe/internal/pkg/response"
	"time"
)

type Repo interface {
	FindByUid(context.Context, string) (*User, error)
	FindById(context.Context, string) (*User, error)
	FindByUsername(context.Context, string) (*User, error)
	Create(context.Context, *User) (*User, error)
	FindByUserEmail(ctx context.Context, name string) (*User, error)
}

type Service struct {
	repo   Repo
	rdb    *redis.Client
	sf     *sonyflake.Sonyflake
	logger *zap.Logger
}

func NewService(repo Repo, rdb *redis.Client, sf *sonyflake.Sonyflake, logger *zap.Logger) *Service {
	return &Service{repo: repo, rdb: rdb, sf: sf, logger: logger}
}

func (s *Service) Register(ctx *gin.Context, param *user.Register) (*User, error) {
	userInstance := &User{
		Name:     param.Username,
		Email:    param.UserEmail,
		Password: param.Password,
	}
	userId, err := s.sf.NextID()
	if err != nil {
		s.logger.Error(err.Error())
		return nil, response.InternalServer("register fail; generate userid failed")
	}
	userInstance.Userid = convertor.ToString(userId)

	u, err := s.repo.Create(ctx, userInstance)
	if err != nil {
		return nil, response.BadRequest("register fail; create user failed; " + err.Error())

	}
	return u, nil
}

// valid password
func isCorrectPassword(password, hash string) bool {
	return cryptor.Sha256(password) != hash
}

func (s *Service) Login(ctx *gin.Context, params *user.Login) (*User, error) {
	query, err := s.repo.FindByUserEmail(ctx, params.UserEmail)
	if err != nil {
		return nil, response.BadRequest("find user by email failed; " + err.Error())
	}
	if isCorrectPassword(params.Password, query.Password) {
		return nil, response.BadRequest("login failed; wrong password")
	}

	return query, nil
}

// UserInfo get user info
func (s *Service) UserInfo(ctx *gin.Context, userUid string) (*User, error) {
	u, err := s.repo.FindByUid(ctx, userUid)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, response.InternalServer("find user by uid failed; " + err.Error())
	}
	return u, nil
}

const GuardName = "app"

type JwtRepo interface {
	JoinBlackList(ctx context.Context, userUid, token string, joinUnix int64, expires time.Duration) error
	GetBlackJoinUnix(ctx context.Context, userUid, token string) (int64, error)
}

type JwtService struct {
	conf        *config.Configuration
	logger      *zap.Logger
	repo        JwtRepo
	userService *Service
	lockBuilder *compo.LockBuilder
}

func NewJwtService(
	conf *config.Configuration,
	logger *zap.Logger,
	repo JwtRepo,
	userService *Service,
	lockBuilder *compo.LockBuilder,
) *JwtService {
	return &JwtService{
		conf:        conf,
		logger:      logger,
		repo:        repo,
		userService: userService,
		lockBuilder: lockBuilder,
	}
}

type TokenOutput struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Uid         string `json:"uid"`
}

type CustomerClaims struct {
	Key string `json:"key,omitempty"`
	jwt.RegisteredClaims
}

type JwtUser interface {
	GetUid() string
}

// CreateToken create token
func (s *JwtService) CreateToken(jwtUser JwtUser) (*TokenOutput, *jwt.Token, error) {
	innerToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomerClaims{
		Key: GuardName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(s.conf.Jwt.JwtTtl))),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * (-1000))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jwtUser.GetUid(),
		},
	})
	token, err := innerToken.SignedString([]byte(s.conf.Jwt.Secret))
	if err != nil {
		return nil, nil, err
	}
	return &TokenOutput{
		Uid:         jwtUser.GetUid(),
		AccessToken: token,
		ExpiresIn:   int64(int(s.conf.Jwt.JwtTtl)),
	}, innerToken, nil
}

// JoinBlackList add user to black
func (s *JwtService) JoinBlackList(ctx *gin.Context, userUid string, token *jwt.Token) error {
	nowUnix := time.Now().Unix()
	// expire user token now
	timer := token.Claims.(*CustomerClaims).ExpiresAt.Sub(time.Now())
	log.Println("add black timer", userUid, nowUnix, timer)
	s.logger.Info("add black: " + userUid)

	if err := s.repo.JoinBlackList(ctx, userUid, token.Raw, nowUnix, timer); err != nil {
		s.logger.Error("failed to join black list", zap.String("userUid", userUid), zap.Error(err))
		return err
	}
	return nil
}

// is in black grace period
func isOutGracePeriod(joinUnix, period int64) bool {
	return time.Now().Unix()-joinUnix <= period
}

// IsInBlackList is user in black
func (s *JwtService) IsInBlackList(ctx *gin.Context, userUid, token string) bool {
	joinUnix, err := s.repo.GetBlackJoinUnix(ctx, userUid, token)
	if err != nil || isOutGracePeriod(joinUnix, s.conf.Jwt.JwtBlacklistGracePeriod) {
		return false
	}
	return true
}

// RefreshToken reset user token
func (s *JwtService) RefreshToken(ctx *gin.Context, userUid string, token *jwt.Token) (*TokenOutput, error) {
	tokenId := token.Claims.(*CustomerClaims).ID
	lock := s.lockBuilder.NewLock(ctx, "refresh_token_lock:"+tokenId, s.conf.Jwt.JwtBlacklistGracePeriod)
	defer lock.Release()
	if !lock.Get() {
		return nil, response.InternalServer("refresh token lock expired")
	}
	info, err := s.userService.UserInfo(ctx, userUid)
	if err != nil {
		s.logger.Error("failed to refresh token", zap.String("userUid", userUid), zap.Error(err))
		return nil, err
	}
	tokenOutput, _, err := s.CreateToken(info)
	if err != nil {
		s.logger.Error("failed to create token", zap.String("userUid", userUid), zap.Error(err))
		return nil, err
	}
	err = s.JoinBlackList(ctx, userUid, token)
	if err != nil {
		s.logger.Error("failed to join black", zap.String("userUid", userUid), zap.Error(err))
		return nil, err
	}
	return tokenOutput, nil
}

package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"matasidik/backend/internal/config"
	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository"
)

type AuthService struct {
	users      repository.UserRepository
	sessions   repository.AuthRepository
	key        paseto.V4SymmetricKey
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type AuthResult struct {
	AccessToken    string       `json:"access_token"`
	AccessExpires  string       `json:"access_expires_at"`
	RefreshToken   string       `json:"-"`
	RefreshExpires string       `json:"-"`
	User           *domain.User `json:"user"`
	ActingAs       *domain.User `json:"acting_as,omitempty"`
}

type AuthContext struct {
	User         *domain.User
	Impersonator *domain.User
}

type impersonatorContextKey struct{}

func WithImpersonator(ctx context.Context, user *domain.User) context.Context {
	if user == nil {
		return ctx
	}
	return context.WithValue(ctx, impersonatorContextKey{}, user)
}

func ImpersonatorFromContext(ctx context.Context) *domain.User {
	user, _ := ctx.Value(impersonatorContextKey{}).(*domain.User)
	return user
}

func NewAuthService(cfg config.Config, users repository.UserRepository, sessions repository.AuthRepository) (*AuthService, error) {
	key, err := authKey(cfg.AuthKeyHex)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(cfg.AuthKeyHex) == "" {
		slog.Warn("APP_AUTH_KEY kosong; memakai kunci acak sementara. Token tidak berlaku setelah restart.")
	}
	return &AuthService{
		users: users, sessions: sessions, key: key,
		accessTTL: 30 * time.Minute, refreshTTL: 7 * 24 * time.Hour,
	}, nil
}

func authKey(hexKey string) (paseto.V4SymmetricKey, error) {
	if strings.TrimSpace(hexKey) == "" {
		return paseto.NewV4SymmetricKey(), nil
	}
	key, err := paseto.V4SymmetricKeyFromHex(strings.TrimSpace(hexKey))
	if err != nil {
		return paseto.V4SymmetricKey{}, fmt.Errorf("APP_AUTH_KEY tidak valid: %w", err)
	}
	return key, nil
}

func (s *AuthService) Login(ctx context.Context, username, password, userAgent string) (*AuthResult, error) {
	auth, err := s.users.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, domain.Unauthorized("username atau password salah")
	}
	if !auth.Aktif {
		return nil, domain.Forbidden("akun tidak aktif")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(auth.PasswordHash), []byte(password)); err != nil {
		return nil, domain.Unauthorized("username atau password salah")
	}
	refresh, session, err := s.newSession(auth.ID, userAgent)
	if err != nil {
		return nil, err
	}
	if err := s.sessions.CreateSession(ctx, session); err != nil {
		return nil, err
	}
	_ = s.sessions.UpdateLastLogin(ctx, auth.ID, now())
	user, err := s.users.GetUser(ctx, auth.ID)
	if err != nil {
		return nil, err
	}
	access, expires, err := s.issue(user, nil)
	if err != nil {
		return nil, err
	}
	return &AuthResult{
		AccessToken: access, AccessExpires: expires, RefreshToken: refresh,
		RefreshExpires: session.ExpiresAt, User: user,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken, userAgent string) (*AuthResult, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, domain.Unauthorized("refresh token tidak ada")
	}
	hash := hashToken(refreshToken)
	session, err := s.sessions.GetSessionByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	if session == nil || session.RevokedAt != nil || session.ExpiresAt < now() {
		return nil, domain.Unauthorized("sesi tidak valid atau kedaluwarsa")
	}
	user, err := s.users.GetUser(ctx, session.UserID)
	if err != nil || !user.Aktif {
		return nil, domain.Unauthorized("sesi tidak valid")
	}
	if err := s.sessions.RevokeSession(ctx, session.ID, now()); err != nil {
		return nil, err
	}
	refresh, newSession, err := s.newSession(user.ID, userAgent)
	if err != nil {
		return nil, err
	}
	if err := s.sessions.CreateSession(ctx, newSession); err != nil {
		return nil, err
	}
	access, expires, err := s.issue(user, nil)
	if err != nil {
		return nil, err
	}
	return &AuthResult{
		AccessToken: access, AccessExpires: expires, RefreshToken: refresh,
		RefreshExpires: newSession.ExpiresAt, User: user,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return nil
	}
	session, err := s.sessions.GetSessionByHash(ctx, hashToken(refreshToken))
	if err != nil || session == nil {
		return err
	}
	return s.sessions.RevokeSession(ctx, session.ID, now())
}

type UpdateProfileInput struct {
	Nama string `json:"nama"`
}

func (s *AuthService) UpdateProfile(ctx context.Context, actor *domain.User, in UpdateProfileInput) (*domain.User, error) {
	nama := strings.TrimSpace(in.Nama)
	if len([]rune(nama)) < 2 || len([]rune(nama)) > 100 {
		return nil, domain.Validation("nama harus 2–100 karakter", nil)
	}
	if actor == nil || actor.ID == "" {
		return nil, domain.Unauthorized("sesi tidak valid")
	}
	if err := s.users.UpdateUserName(ctx, actor.ID, nama); err != nil {
		return nil, err
	}
	return s.users.GetUser(ctx, actor.ID)
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (s *AuthService) ChangePassword(ctx context.Context, actor *domain.User, in ChangePasswordInput) error {
	if actor == nil || actor.ID == "" || strings.TrimSpace(actor.Username) == "" {
		return domain.Unauthorized("sesi tidak valid")
	}
	if len(in.NewPassword) < 8 {
		return domain.Validation("password baru minimal 8 karakter", nil)
	}
	if in.CurrentPassword == in.NewPassword {
		return domain.Validation("password baru harus berbeda dari password saat ini", nil)
	}
	auth, err := s.users.GetUserByUsername(ctx, actor.Username)
	if err != nil {
		return domain.Unauthorized("sesi tidak valid")
	}
	if bcrypt.CompareHashAndPassword([]byte(auth.PasswordHash), []byte(in.CurrentPassword)) != nil {
		return domain.Validation("password saat ini salah", nil)
	}
	hash, err := HashPassword(in.NewPassword)
	if err != nil {
		return err
	}
	if err := s.users.UpdateUserPassword(ctx, actor.ID, hash); err != nil {
		return err
	}
	return s.sessions.RevokeSessionsForUser(ctx, actor.ID, now())
}

func (s *AuthService) ActingAs(ctx context.Context, admin *domain.User, targetID string) (*AuthResult, error) {
	if admin.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh actingAs")
	}
	target, err := s.users.GetUser(ctx, targetID)
	if err != nil {
		return nil, err
	}
	if target.Role == domain.RoleAdmin {
		return nil, domain.Forbidden("tidak boleh actingAs sesama admin")
	}
	if !target.Aktif {
		return nil, domain.Forbidden("akun target tidak aktif")
	}
	access, expires, err := s.issue(target, admin)
	if err != nil {
		return nil, err
	}
	return &AuthResult{AccessToken: access, AccessExpires: expires, User: target, ActingAs: admin}, nil
}

func (s *AuthService) Parse(ctx context.Context, raw string) (*AuthContext, error) {
	parser := paseto.NewParserForValidNow()
	token, err := parser.ParseV4Local(s.key, strings.TrimSpace(raw), nil)
	if err != nil {
		return nil, domain.Unauthorized("token tidak valid")
	}
	sub, err := token.GetSubject()
	if err != nil || sub == "" {
		return nil, domain.Unauthorized("token tidak valid")
	}
	user, err := s.users.GetUser(ctx, sub)
	if err != nil || !user.Aktif {
		return nil, domain.Unauthorized("user tidak aktif")
	}
	ac := &AuthContext{User: user}
	if act, err := token.GetString("act"); err == nil && act != "" {
		impersonator, err := s.users.GetUser(ctx, act)
		if err != nil || impersonator.Role != domain.RoleAdmin || !impersonator.Aktif {
			return nil, domain.Unauthorized("actingAs tidak valid")
		}
		ac.Impersonator = impersonator
	}
	return ac, nil
}

func (s *AuthService) issue(user *domain.User, impersonator *domain.User) (string, string, error) {
	nowTime := time.Now().UTC()
	exp := nowTime.Add(s.accessTTL)
	token := paseto.NewToken()
	token.SetSubject(user.ID)
	token.SetIssuer("matasidik")
	token.SetAudience("matasidik-web")
	token.SetIssuedAt(nowTime)
	token.SetNotBefore(nowTime)
	token.SetExpiration(exp)
	token.SetJti(uuid.NewString())
	token.SetString("role", string(user.Role))
	if user.FaskesID != nil {
		token.SetString("faskes_id", *user.FaskesID)
	}
	if user.BidangID != nil {
		token.SetString("bidang_id", *user.BidangID)
	}
	if impersonator != nil {
		token.SetString("act", impersonator.ID)
	}
	raw := token.V4Encrypt(s.key, nil)
	return raw, exp.Format(time.RFC3339), nil
}

func (s *AuthService) newSession(userID, userAgent string) (string, *domain.AuthSession, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", nil, err
	}
	refresh := hex.EncodeToString(buf)
	nowTime := time.Now().UTC()
	session := &domain.AuthSession{
		ID: uuid.NewString(), UserID: userID, RefreshHash: hashToken(refresh),
		UserAgent: strPtr(userAgent), CreatedAt: nowTime.Format(time.RFC3339),
		ExpiresAt: nowTime.Add(s.refreshTTL).Format(time.RFC3339),
	}
	return refresh, session, nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(b), err
}

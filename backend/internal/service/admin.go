package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository"
)

type AdminService struct {
	repo repository.Repository
}

func NewAdminService(repo repository.Repository) *AdminService { return &AdminService{repo: repo} }

func (s *AdminService) Users(ctx context.Context) ([]domain.User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *AdminService) Positions(ctx context.Context) ([]domain.Position, error) {
	return s.repo.ListPositions(ctx)
}

func (s *AdminService) Faskes(ctx context.Context) ([]domain.Faskes, error) {
	return s.repo.ListFaskes(ctx)
}

func (s *AdminService) effectiveRole(ctx context.Context, positionID string, fallback domain.Role) (domain.Role, error) {
	if strings.TrimSpace(positionID) == "" {
		if fallback == "" {
			return domain.RolePKMOperator, nil
		}
		return fallback, nil
	}
	positions, err := s.repo.ListPositions(ctx)
	if err != nil {
		return "", err
	}
	for _, p := range positions {
		if p.ID == positionID {
			return p.Role, nil
		}
	}
	return "", domain.Validation("position_id tidak ditemukan", nil)
}

type CreateUserInput struct {
	Username   string `json:"username"`
	Nama       string `json:"nama"`
	Password   string `json:"password"`
	PositionID string `json:"position_id"`
	FaskesID   string `json:"faskes_id"`
	BidangID   string `json:"bidang_id"`
}

func (s *AdminService) CreateUser(ctx context.Context, in CreateUserInput, actor *domain.User) (*domain.User, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh membuat user")
	}
	if strings.TrimSpace(in.Username) == "" || strings.TrimSpace(in.Nama) == "" || len(in.Password) < 8 {
		return nil, domain.Validation("username, nama, dan password minimal 8 karakter wajib diisi", nil)
	}
	role, err := s.effectiveRole(ctx, in.PositionID, domain.RolePKMOperator)
	if err != nil {
		return nil, err
	}
	hash, err := HashPassword(in.Password)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		ID: uuid.NewString(), Username: strings.TrimSpace(in.Username), Nama: strings.TrimSpace(in.Nama),
		Role: role, PositionID: strPtr(in.PositionID), FaskesID: strPtr(in.FaskesID),
		BidangID: strPtr(in.BidangID), Aktif: true, CreatedAt: now(),
	}
	if err := s.repo.CreateUser(ctx, user, hash); err != nil {
		return nil, fmt.Errorf("buat user: %w", err)
	}
	return s.repo.GetUser(ctx, user.ID)
}

type UpdateUserInput struct {
	Username   string `json:"username"`
	Nama       string `json:"nama"`
	PositionID string `json:"position_id"`
	FaskesID   string `json:"faskes_id"`
	BidangID   string `json:"bidang_id"`
	Aktif      *bool  `json:"aktif"`
}

func (s *AdminService) UpdateUser(ctx context.Context, id string, in UpdateUserInput, actor *domain.User) (*domain.User, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh mengubah user")
	}
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.Nama) != "" {
		user.Nama = strings.TrimSpace(in.Nama)
	}
	if strings.TrimSpace(in.Username) != "" {
		user.Username = strings.TrimSpace(in.Username)
	}
	user.PositionID = strPtr(in.PositionID)
	user.FaskesID = strPtr(in.FaskesID)
	user.BidangID = strPtr(in.BidangID)
	if in.PositionID != "" {
		role, err := s.effectiveRole(ctx, in.PositionID, user.Role)
		if err != nil {
			return nil, err
		}
		user.Role = role
	}
	if in.Aktif != nil {
		user.Aktif = *in.Aktif
	}
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return s.repo.GetUser(ctx, id)
}

func (s *AdminService) ResetPassword(ctx context.Context, id, password string, actor *domain.User) error {
	if actor.Role != domain.RoleAdmin {
		return domain.Forbidden("hanya admin yang boleh reset password")
	}
	if len(password) < 8 {
		return domain.Validation("password minimal 8 karakter", nil)
	}
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	return s.repo.UpdateUserPassword(ctx, id, hash)
}

type FaskesInput struct {
	Kode      string `json:"kode"`
	Nama      string `json:"nama"`
	WilayahID string `json:"wilayah_id"`
	Jenis     string `json:"jenis"`
	Aktif     *bool  `json:"aktif"`
}

func (s *AdminService) CreateFaskes(ctx context.Context, in FaskesInput, actor *domain.User) (*domain.Faskes, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh membuat faskes")
	}
	if strings.TrimSpace(in.Kode) == "" || strings.TrimSpace(in.Nama) == "" {
		return nil, domain.Validation("kode dan nama faskes wajib diisi", nil)
	}
	if in.Jenis == "" {
		in.Jenis = "puskesmas"
	}
	f := &domain.Faskes{
		ID: uuid.NewString(), Kode: strings.TrimSpace(in.Kode), Nama: strings.TrimSpace(in.Nama),
		WilayahID: strPtr(in.WilayahID), Jenis: in.Jenis, Aktif: true,
	}
	if err := s.repo.CreateFaskes(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *AdminService) UpdateFaskes(ctx context.Context, id string, in FaskesInput, actor *domain.User) (*domain.Faskes, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh mengubah faskes")
	}
	f, err := s.repo.GetFaskes(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Kode != "" {
		f.Kode = in.Kode
	}
	if in.Nama != "" {
		f.Nama = in.Nama
	}
	f.WilayahID = strPtr(in.WilayahID)
	if in.Jenis != "" {
		f.Jenis = in.Jenis
	}
	if in.Aktif != nil {
		f.Aktif = *in.Aktif
	}
	if err := s.repo.UpdateFaskes(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

type PositionInput struct {
	Kode     string `json:"kode"`
	Nama     string `json:"nama"`
	Role     string `json:"role"`
	Level    int    `json:"level"`
	BidangID string `json:"bidang_id"`
	FaskesID string `json:"faskes_id"`
	Aktif    *bool  `json:"aktif"`
}

func (s *AdminService) CreatePosition(ctx context.Context, in PositionInput, actor *domain.User) (*domain.Position, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh membuat posisi")
	}
	if in.Kode == "" || in.Nama == "" || in.Role == "" {
		return nil, domain.Validation("kode, nama, dan role posisi wajib diisi", nil)
	}
	if in.Level == 0 {
		in.Level = 1
	}
	p := &domain.Position{
		ID: uuid.NewString(), Kode: in.Kode, Nama: in.Nama, Role: domain.Role(in.Role),
		Level: in.Level, BidangID: strPtr(in.BidangID), FaskesID: strPtr(in.FaskesID),
		Aktif: true, CreatedAt: now(),
	}
	if err := s.repo.CreatePosition(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *AdminService) UpdatePosition(ctx context.Context, id string, in PositionInput, actor *domain.User) (*domain.Position, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh mengubah posisi")
	}
	positions, err := s.repo.ListPositions(ctx)
	if err != nil {
		return nil, err
	}
	var p *domain.Position
	for i := range positions {
		if positions[i].ID == id {
			p = &positions[i]
			break
		}
	}
	if p == nil {
		return nil, domain.NotFound("posisi")
	}
	if in.Kode != "" {
		p.Kode = in.Kode
	}
	if in.Nama != "" {
		p.Nama = in.Nama
	}
	if in.Role != "" {
		p.Role = domain.Role(in.Role)
	}
	if in.Level > 0 {
		p.Level = in.Level
	}
	p.BidangID = strPtr(in.BidangID)
	p.FaskesID = strPtr(in.FaskesID)
	if in.Aktif != nil {
		p.Aktif = *in.Aktif
	}
	if err := s.repo.UpdatePosition(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

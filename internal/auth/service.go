// internal/auth/service.go
package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

// ── Register ────────────────────────────────────────────────

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("name, email and password are required")
	}
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	if req.Role == "" {
		req.Role = "both"
	}

	// check email not already taken
	var exists bool
	err := s.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, req.Email,
	).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// insert user
	var id, name, email, role string
	err = s.db.QueryRow(ctx,
		`INSERT INTO users (name, email, password_hash, phone, role)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, name, email, role`,
		req.Name, req.Email, string(hash), req.Phone, req.Role,
	).Scan(&id, &name, &email, &role)
	if err != nil {
		return nil, err
	}

	token, err := generateToken(id, email, role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  UserPublic{ID: id, Name: name, Email: email, Role: role},
	}, nil
}

// ── Login ───────────────────────────────────────────────────

func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	var id, name, email, role, hash string
	err := s.db.QueryRow(ctx,
		`SELECT id, name, email, role, password_hash
		 FROM users WHERE email = $1`, req.Email,
	).Scan(&id, &name, &email, &role, &hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("invalid email or password")
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := generateToken(id, email, role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  UserPublic{ID: id, Name: name, Email: email, Role: role},
	}, nil
}

// ── JWT helpers ─────────────────────────────────────────────

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims // This is a Go feature called Embedding (or Composition). It allows you to include the fields and methods of one struct (RegisteredClaims) directly within another struct (Claims) without needing to explicitly define them. This is a common pattern in Go for code reuse and to create more complex data structures by combining simpler ones.
}

func generateToken(userID, email, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenStr string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username      string             `bson:"username" json:"username"`
	Email         string             `bson:"email" json:"email"`
	PasswordHash  string             `bson:"password_hash" json:"-"`
	DisplayName   string             `bson:"display_name" json:"display_name"`
	Bio           string             `bson:"bio" json:"bio"`
	AvatarURL     string             `bson:"avatar_url" json:"avatar_url"`
	FollowerCount int                `bson:"follower_count" json:"follower_count"`
	FollowingCount int               `bson:"following_count" json:"following_count"`
	VideoCount    int                `bson:"video_count" json:"video_count"`
	TotalLikes    int                `bson:"total_likes" json:"total_likes"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=30"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

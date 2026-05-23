package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"abstract-api/internal/db"
)

var ErrNotFound = errors.New("not found")

type Profile struct {
	UID             string `json:"uid"`
	Name            string `json:"name"`
	Gender          string `json:"gender,omitempty"`
	ShellBalance    int64  `json:"shellBalance"`
	AvatarIdx       *int   `json:"avatarIdx,omitempty"`
	ActiveFrame     string `json:"activeFrame,omitempty"`
	ActiveAccessory string `json:"activeAccessory,omitempty"`
}

type Trail struct {
	ID          int    `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	OrderIndex  int    `json:"orderIndex"`
	Published   bool   `json:"published"`
}

type Mission struct {
	ID           int    `json:"id"`
	TrailID      int    `json:"trailId"`
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	ContentMd    string `json:"contentMd,omitempty"`
	RewardShells int    `json:"rewardShells"`
	OrderIndex   int    `json:"orderIndex"`
	Published    bool   `json:"published"`
}

func EnsureUser(ctx context.Context, uid string) error {
	if db.Pool == nil {
		return fmt.Errorf("db not initialized")
	}
	_, err := db.Pool.Exec(ctx, `INSERT INTO users (uid) VALUES ($1) ON CONFLICT (uid) DO NOTHING`, uid)
	return err
}

func GetProfile(ctx context.Context, uid string) (*Profile, error) {
	if db.Pool == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	if err := EnsureUser(ctx, uid); err != nil {
		return nil, err
	}

	row := db.Pool.QueryRow(ctx, `
		SELECT uid, COALESCE(name, ''), COALESCE(gender, ''), shell_balance, avatar_idx, COALESCE(active_frame, ''), COALESCE(active_accessory, '')
		FROM users
		WHERE uid = $1
	`, uid)

	var profile Profile
	var gender, activeFrame, activeAccessory string
	var avatarIdx sql.NullInt32
	if err := row.Scan(&profile.UID, &profile.Name, &gender, &profile.ShellBalance, &avatarIdx, &activeFrame, &activeAccessory); err != nil {
		return nil, err
	}
	profile.Gender = gender
	if avatarIdx.Valid {
		idx := int(avatarIdx.Int32)
		profile.AvatarIdx = &idx
	}
	profile.ActiveFrame = activeFrame
	profile.ActiveAccessory = activeAccessory
	return &profile, nil
}

type ProfileUpdate struct {
	Name            string `json:"name"`
	Gender          string `json:"gender"`
	AvatarIdx       *int   `json:"avatarIdx"`
	ActiveFrame     string `json:"activeFrame"`
	ActiveAccessory string `json:"activeAccessory"`
}

func UpdateProfile(ctx context.Context, uid string, input ProfileUpdate) (*Profile, error) {
	if db.Pool == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	if err := EnsureUser(ctx, uid); err != nil {
		return nil, err
	}

	_, err := db.Pool.Exec(ctx, `
		UPDATE users
		SET name = COALESCE(NULLIF($2, ''), name),
		    gender = COALESCE(NULLIF($3, ''), gender),
		    avatar_idx = COALESCE($4, avatar_idx),
		    active_frame = COALESCE(NULLIF($5, ''), active_frame),
		    active_accessory = COALESCE(NULLIF($6, ''), active_accessory),
		    updated_at = now()
		WHERE uid = $1
	`, uid, input.Name, input.Gender, input.AvatarIdx, input.ActiveFrame, input.ActiveAccessory)
	if err != nil {
		return nil, err
	}
	return GetProfile(ctx, uid)
}

func UpdateAvatar(ctx context.Context, uid string, avatarIdx int) (*Profile, error) {
	if db.Pool == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	if err := EnsureUser(ctx, uid); err != nil {
		return nil, err
	}
	_, err := db.Pool.Exec(ctx, `
		UPDATE users
		SET avatar_idx = $2,
		    updated_at = now()
		WHERE uid = $1
	`, uid, avatarIdx)
	if err != nil {
		return nil, err
	}
	return GetProfile(ctx, uid)
}

func ListTrails(ctx context.Context) ([]Trail, error) {
	if db.Pool == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	rows, err := db.Pool.Query(ctx, `
		SELECT id, slug, title, COALESCE(description, ''), COALESCE(order_index, 0), published
		FROM trails
		WHERE published = true OR published IS NULL
		ORDER BY COALESCE(order_index, id)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trails := make([]Trail, 0)
	for rows.Next() {
		var item Trail
		if err := rows.Scan(&item.ID, &item.Slug, &item.Title, &item.Description, &item.OrderIndex, &item.Published); err != nil {
			return nil, err
		}
		trails = append(trails, item)
	}
	return trails, rows.Err()
}

func GetTrail(ctx context.Context, trailID string) (*Trail, error) {
	if db.Pool == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	_, err := strconv.Atoi(trailID)
	if err != nil {
		return nil, fmt.Errorf("trailId must be numeric")
	}

	row := db.Pool.QueryRow(ctx, `
		SELECT id, slug, title, COALESCE(description, ''), COALESCE(order_index, 0), published
		FROM trails
		WHERE id = $1
	`, trailID)

	var item Trail
	if err := row.Scan(&item.ID, &item.Slug, &item.Title, &item.Description, &item.OrderIndex, &item.Published); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

func ListTrailMissions(ctx context.Context, trailID string) ([]Mission, error) {
	if db.Pool == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	_, err := strconv.Atoi(trailID)
	if err != nil {
		return nil, fmt.Errorf("trailId must be numeric")
	}

	rows, err := db.Pool.Query(ctx, `
		SELECT id, trail_id, slug, title, COALESCE(content_md, ''), COALESCE(reward_shells, 0), COALESCE(order_index, 0), published
		FROM missions
		WHERE trail_id = $1 AND (published = true OR published IS NULL)
		ORDER BY COALESCE(order_index, id)
	`, trailID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	missions := make([]Mission, 0)
	for rows.Next() {
		var item Mission
		if err := rows.Scan(&item.ID, &item.TrailID, &item.Slug, &item.Title, &item.ContentMd, &item.RewardShells, &item.OrderIndex, &item.Published); err != nil {
			return nil, err
		}
		missions = append(missions, item)
	}
	return missions, rows.Err()
}

func GetMission(ctx context.Context, missionID string) (*Mission, error) {
	if db.Pool == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	_, err := strconv.Atoi(missionID)
	if err != nil {
		return nil, fmt.Errorf("missionId must be numeric")
	}

	row := db.Pool.QueryRow(ctx, `
		SELECT id, trail_id, slug, title, COALESCE(content_md, ''), COALESCE(reward_shells, 0), COALESCE(order_index, 0), published
		FROM missions
		WHERE id = $1
	`, missionID)

	var item Mission
	if err := row.Scan(&item.ID, &item.TrailID, &item.Slug, &item.Title, &item.ContentMd, &item.RewardShells, &item.OrderIndex, &item.Published); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

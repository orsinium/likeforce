package likeforce

import (
	"fmt"

	"github.com/go-redis/redis"
)

// Users is a set of operations to store user stat in Redis
type Users struct {
	client *redis.Client
}

func makeKeyPosts(chat int64, user int) string {
	return fmt.Sprintf("likes:user:posts:%d:%d", chat, user)
}

func makeKeyRating(chat int64, user int) string {
	return fmt.Sprintf("likes:user:rating:%d:%d", chat, user)
}

func makeKeyName(user int) string {
	return fmt.Sprintf("likes:user:login:%d", user)
}

// AddPost to increment posts count for user
func (storage *Users) AddPost(chat int64, user int) error {
	return storage.client.Incr(makeKeyPosts(chat, user)).Err()
}

// AddRating to increment rating for user
func (storage *Users) AddRating(chat int64, user int) error {
	return storage.client.Incr(makeKeyRating(chat, user)).Err()
}

// AddName to save username
func (storage *Users) AddName(user int, name string) error {
	return storage.client.Set(makeKeyName(user), name, 0).Err()
}

// GetName to get username
func (storage *Users) GetName(user int) (string, error) {
	return storage.client.Get(makeKeyName(user)).Result()
}

// RemoveRating to decrement rating for user
func (storage *Users) RemoveRating(chat int64, user int) error {
	return storage.client.Decr(makeKeyRating(chat, user)).Err()
}

// PostsCount to get posts count for user
func (storage *Users) PostsCount(chat int64, user int) (int, error) {
	key := makeKeyPosts(chat, user)
	keysCount, err := storage.client.Exists(key).Result()
	if err != nil {
		return 0, err
	}
	if keysCount == 0 {
		return 0, nil
	}
	return storage.client.Get(key).Int()
}

// RatingCount to get rating for user
func (storage *Users) RatingCount(chat int64, user int) (int, error) {
	key := makeKeyRating(chat, user)
	keysCount, err := storage.client.Exists(key).Result()
	if err != nil {
		return 0, err
	}
	if keysCount == 0 {
		return 0, nil
	}
	return storage.client.Get(key).Int()
}

// ByteCount to make human-readable rating
func ByteCount(count int) string {
	const unit = 1000
	if count < unit {
		return fmt.Sprintf("%d", count)
	}
	div, exp := int64(unit), 0
	for n := count / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(count)/float64(div), "kMGTPE"[exp])
}

// Stat to get human-readable message with user stat
func (storage *Users) Stat(chat int64, user int) (string, error) {
	posts, err := storage.PostsCount(chat, user)
	if err != nil {
		return "", err
	}
	rating, err := storage.RatingCount(chat, user)
	if err != nil {
		return "", err
	}
	if posts == 0 {
		return "First blood!", nil
	}
	const tmpl = "user stat:\nposts: %s\nrating: %s"
	return fmt.Sprintf(tmpl, ByteCount(posts), ByteCount(rating)), nil
}

package model

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

const (
	// 五分钟过期
	expireTime = 5 * time.Minute
)

type Captcha struct {
	Key     string `json:"-"`
	Content string
	Purpose int32
}

func (c *Captcha) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
func (c *Captcha) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

// 设置验证码
func (p *Provider) SetCaptcha(c Captcha) error {
	return p.cache.Set(c.Key, &c, expireTime).Err()
}

func (p *Provider) Verify(c Captcha) (bool, error) {
	cmd := p.cache.Get(c.Key)
	if err := cmd.Err(); err != nil {
		return false, err
	}

	var val Captcha
	if err := cmd.Scan(&val); err != nil && err != redis.Nil {
		return false, err
	}
	if c.Content != val.Content ||
		c.Purpose != val.Purpose {
		return false, nil
	}

	p.cache.Del(c.Key)
	return true, nil
}

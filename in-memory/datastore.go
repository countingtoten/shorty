package memory

import (
	"errors"

	"github.com/countingtoten/shorty"
	"github.com/countingtoten/shorty/rand"
)

var (
	ErrShortURLNotFound = errors.New("missing short url")
	ErrUserNotFound     = errors.New("missing user")
)

type Datastore struct {
	shorty.Config
	UserData   map[shorty.UserID]*shorty.User
	ShortCodes map[shorty.ShortCode]*shorty.URL
}

func NewDatastore(cfg shorty.Config) *Datastore {
	return &Datastore{
		Config:     cfg,
		UserData:   map[shorty.UserID]*shorty.User{},
		ShortCodes: map[shorty.ShortCode]*shorty.URL{},
	}
}

func (d *Datastore) CreateShortURL(id shorty.UserID, url shorty.LongURL) (shorty.ShortURL, error) {
	user := d.GetUser(id)

	shortCode := d.NewShortCode()

	newURL := &shorty.URL{
		ShortCode: shortCode,
		LongURL:   url,
	}

	user.URLs[url] = newURL

	d.ShortCodes[shortCode] = newURL

	shortURL := d.Config.BaseURL + shortCode

	return shortURL, nil
}

func (d *Datastore) GetLongURL(url shorty.ShortCode) (shorty.LongURL, error) {
	data, ok := d.ShortCodes[url]
	if !ok {
		return "", nil
	}

	return data.LongURL, nil
}

func (d *Datastore) GetUser(id shorty.UserID) *shorty.User {
	user, ok := d.UserData[id]

	if !ok {
		// TODO: Add user authentication
		// Until then just add any new users
		user = &shorty.User{
			ID:   id,
			URLs: map[shorty.LongURL]*shorty.URL{},
		}

		d.UserData[id] = user
	}

	return user
}

func (d *Datastore) NewShortCode() shorty.ShortCode {
	for {
		code := rand.AlphanumericString(d.Config.ShortCodeLength)

		_, ok := d.ShortCodes[code]

		if !ok {
			return code
		}
	}
}

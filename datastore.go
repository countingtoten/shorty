package shorty

//go:generate mockery -name=Datastore -case=underscore

type Datastore interface {
	// GetLongURL returns an empty string if the short url does not exist and it
	// will return an error if there is an error accessing the datastore
	GetLongURL(url ShortCode) (LongURL, error)
	CreateShortURL(id UserID, url LongURL) (ShortURL, error)
}

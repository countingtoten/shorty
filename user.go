package shorty

type UserID = int64

type User struct {
	ID   UserID
	URLs map[LongURL]*URL
}

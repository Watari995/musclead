package healthsyncdomain

import "github.com/Watari995/musclead/internal/valueobject"

type StateSigner interface {
	Sign(userID valueobject.UserID, redirectURL string) (string, error)
	Verify(state string) (userID valueobject.UserID, redirectURL string, err error)
}

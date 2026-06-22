package healthsyncdomain

import "github.com/Watari995/musclead/internal/valueobject"

type StateSigner interface {
	Sign(userID valueobject.UserID) (string, error)
	Verify(state string) (valueobject.UserID, error)
}

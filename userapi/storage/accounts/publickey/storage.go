package publickey

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/blake2b"

	"github.com/matrix-org/dendrite/clientapi/auth/authtypes"
	"github.com/matrix-org/dendrite/internal/sqlutil"
	"github.com/matrix-org/dendrite/userapi/api"
	"github.com/matrix-org/gomatrixserverlib"
)

// Database represents a read only pseudo database for a pubkey authentication
type Database struct {
	ServerName gomatrixserverlib.ServerName
}

func verifySignedTimeWindow(publicKey ed25519.PublicKey, tw int64, sig []byte) bool {
	enquiry := fmt.Sprintf("login:%d", tw)
	digest := blake2b.Sum256([]byte(enquiry))
	return ed25519.Verify(publicKey, digest[:], sig)
}

func (d *Database) GetAccountByPassword(ctx context.Context, localpart, plaintextPassword string) (*api.Account, error) {
	p := strings.Split(plaintextPassword, ":")
	if len(p) != 3 || p[0] != "ed" {
		return nil, errors.New("error parsing public key credentials")
	}

	sig, err := hex.DecodeString(p[1])
	if err != nil {
		return nil, err
	}
	pub, err := hex.DecodeString(p[2])
	if err != nil {
		return nil, err
	}
	clientHash, err := hex.DecodeString(localpart)
	if err != nil {
		return nil, err
	}
	pubHash := blake2b.Sum256(pub)
	if !bytes.Equal(clientHash, pubHash[:]) {
		return nil, errors.New("public key hash doesn't match public key")
	}

	tw := time.Now().Unix() / (5 * 60)
	if !verifySignedTimeWindow(ed25519.PublicKey(pub), tw, sig) && !verifySignedTimeWindow(ed25519.PublicKey(pub), tw-1, sig) {
		return nil, errors.New("invalid signature")
	}

	return &api.Account{
		UserID:     fmt.Sprintf("@%s:%s", localpart, string(d.ServerName)),
		Localpart:  localpart,
		ServerName: d.ServerName,
	}, nil
}

func (d *Database) GetAccountByLocalpart(ctx context.Context, localpart string) (*api.Account, error) {
	return &api.Account{
		UserID:     fmt.Sprintf("@%s:%s", localpart, string(d.ServerName)),
		Localpart:  localpart,
		ServerName: d.ServerName,
	}, nil
}

func (d *Database) GetProfileByLocalpart(ctx context.Context, localpart string) (*authtypes.Profile, error) {
	return &authtypes.Profile{Localpart: localpart}, nil
}

// functions below are just stubs

func (d *Database) SetPassword(ctx context.Context, localpart string, plaintextPassword string) error {
	return errors.New("can't set password in a public key only mode")
}

func (d *Database) SetAvatarURL(ctx context.Context, localpart string, avatarURL string) error {
	return errors.New("can't set avatar in a public key only mode")
}

func (d *Database) SetDisplayName(ctx context.Context, localpart string, displayName string) error {
	return errors.New("can't set name in a public key only mode")
}

func (d *Database) CreateAccount(ctx context.Context, localpart, plaintextPassword, appserviceID string) (*api.Account, error) {
	return nil, errors.New("can't create acount in a public key only mode")
}

func (d *Database) CreateGuestAccount(ctx context.Context) (*api.Account, error) {
	return nil, errors.New("can't create acount in a public key only mode")
}

func (d *Database) SaveAccountData(ctx context.Context, localpart, roomID, dataType string, content json.RawMessage) error {
	return errors.New("can't set saved data in a public key only mode")
}

func (d *Database) GetAccountData(ctx context.Context, localpart string) (global map[string]json.RawMessage, rooms map[string]map[string]json.RawMessage, err error) {
	return nil, nil, nil
}

func (d *Database) GetAccountDataByType(ctx context.Context, localpart, roomID, dataType string) (data json.RawMessage, err error) {
	return nil, nil
}

func (d *Database) GetNewNumericLocalpart(ctx context.Context) (int64, error) {
	return 0, errors.New("can't generate numeric user ID in a public key only mode")
}

func (d *Database) SaveThreePIDAssociation(ctx context.Context, threepid, localpart, medium string) (err error) {
	return errors.New("can't save 3PID association in a public key only mode")
}

func (d *Database) RemoveThreePIDAssociation(ctx context.Context, threepid string, medium string) (err error) {
	return errors.New("can't save 3PID association in a public key only mode")
}

func (d *Database) GetLocalpartForThreePID(ctx context.Context, threepid string, medium string) (localpart string, err error) {
	return "", nil
}

func (d *Database) GetThreePIDsForLocalpart(ctx context.Context, localpart string) (threepids []authtypes.ThreePID, err error) {
	return []authtypes.ThreePID{}, nil
}

func (d *Database) CheckAccountAvailability(ctx context.Context, localpart string) (bool, error) {
	return true, nil
}

func (d *Database) SearchProfiles(ctx context.Context, searchString string, limit int) ([]authtypes.Profile, error) {
	return []authtypes.Profile{}, nil
}

func (d *Database) DeactivateAccount(ctx context.Context, localpart string) (err error) {
	return errors.New("can't deactivate account in a public key only mode")
}

func (d *Database) CreateOpenIDToken(ctx context.Context, token, localpart string) (exp int64, err error) {
	return 0, errors.New("can't create OpenID token in a public key only mode")
}

func (d *Database) GetOpenIDTokenAttributes(ctx context.Context, token string) (*api.OpenIDTokenAttributes, error) {
	return nil, nil
}

func (d *Database) PartitionOffsets(ctx context.Context, topic string) ([]sqlutil.PartitionOffset, error) {
	return []sqlutil.PartitionOffset{}, nil
}

func (d *Database) SetPartitionOffset(ctx context.Context, topic string, partition int32, offset int64) error {
	return nil
}

package user

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"time"
)

type UserModel struct {
	Driver *neo4j.DriverWithContext
}

// IsUserBlocked Method to check if the relation between current user
// and user to check is block relation
func (u *UserModel) IsUserBlocked(username, userToCheck string) (bool, error) {
	return false, nil
}

// IsUserFriend Method to check if the relation between current user and
// user to check has friend relation with ACCEPTED status
func (u *UserModel) IsUserFriend(username, userToCheck string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := neo4j.ExecuteQuery(
		ctx,
		*u.Driver,
		"RETURN EXISTS {MATCH (:User {username: $username})-[r:IS_FRIEND {status: \"ACCEPTED\"}]->(:User {username: $userToCheck})} AS isAccepted",
		map[string]any{
			"username":    username,
			"userToCheck": userToCheck,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)

	if err != nil {
		log.Printf("Error getting user relationship status: %v", err)
		return false, err
	}

	if len(result.Records) > 0 {
		record := result.Records[0]
		value, _ := record.Get("isAccepted")
		isFriend, _ := value.(bool)

		log.Printf("%q is friend with %q: %v", username, userToCheck, isFriend)
		return isFriend, nil
	}

	return false, nil
}

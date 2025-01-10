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
	query := `
		RETURN EXISTS {
			MATCH (:User {username: $username})-[r:IS_FRIEND {status: "ACCEPTED"}]->(:User {username: $userToCheck})
		} AS isAccepted
	`
	params := map[string]any{
		"username":    username,
		"userToCheck": userToCheck,
	}

	result, err := neo4j.ExecuteQuery(
		ctx,
		*u.Driver,
		query,
		params,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)

	if err != nil {
		log.Printf("Error getting user relationship status: %v\n", err)
		return false, err
	}

	if len(result.Records) > 0 {
		record := result.Records[0]
		value, _ := record.Get("isAccepted")
		isFriend, _ := value.(bool)

		log.Printf("%q is friend with %q: %v\n", username, userToCheck, isFriend)
		return isFriend, nil
	}

	return false, nil
}

func (u *UserModel) getUserDetails(username string, usersToCheck []string) ([]*UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
		MATCH (this:User)
		WHERE this.username IN $users
		CALL {
			WITH this
			MATCH (this)-[this0:IS_FRIEND]-(this1:User)
			WHERE this1.username = $currentUser
			WITH collect({ node: this1, relationship: this0 }) AS edges
			WITH edges
			CALL {
				WITH edges
				UNWIND edges AS edge
				WITH edge.node AS this1, edge.relationship AS this0
				RETURN collect({ requestedBy: this0.requestedBy, addedOn: apoc.date.convertFormat(toString(this0.addedOn), "iso_zoned_date_time", "iso_offset_date_time"), status: this0.status }) AS var2
			}
			RETURN var2
		}
		RETURN this { .username, .name, .id, .profilePicture, friendsConnection: var2 } AS user
	`
	params := map[string]any{
		"users":       usersToCheck,
		"currentUser": username,
	}
	result, err := neo4j.ExecuteQuery(
		ctx,
		*u.Driver,
		query,
		params,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"),
	)
	if err != nil {
		log.Printf("Error getting user info: %v\n", err)
		return nil, err
	}

	users := make([]*UserInfo, 0)
	for _, record := range result.Records {
		userMap := record.AsMap()["user"].(map[string]any)

		friendConnection := userMap["friendsConnection"].([]any)
		var relationInfo *UserRelationInfo
		if len(friendConnection) > 0 {
			info := friendConnection[0].(map[string]any)
			relationInfo = &UserRelationInfo{
				AddedOn:     info["addedOn"].(string),
				RequestedBy: info["requestedBy"].(string),
				Status:      info["status"].(string),
			}
		}
		user := UserInfo{
			Username:       userMap["username"].(string),
			UserId:         userMap["id"].(string),
			Name:           userMap["name"].(string),
			ProfilePicture: userMap["profilePicture"].(string),
			RelationInfo:   relationInfo,
		}

		users = append(users, &user)
	}

	return users, nil
}
package service

import (
	"context"
	"fmt"

	DB "github.com/MRD1920/Notification-System/db"
	db "github.com/MRD1920/Notification-System/db/sqlc"
	model "github.com/MRD1920/Notification-System/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func AddUserToDB(user model.User) error {
	//Add the user to the database
	queries := db.New(DB.Pool)

	dbUser, err := convertUserToDBUser(user)
	if err != nil {
		return err
	}

	err = queries.CreateUser(context.Background(), dbUser)
	if err != nil {
		return err
	}
	return nil

}

func convertUserToDBUser(user model.User) (db.CreateUserParams, error) {

	return db.CreateUserParams{
		ID:                      pgtype.UUID{Bytes: user.Id, Valid: true},
		Username:                user.Username,
		Email:                   user.Email,
		Phone:                   pgtype.Text{String: user.Phone, Valid: true},
		Status:                  user.Status,
		PreferenceLowChannel:    user.Preference.Low,
		PreferenceMediumChannel: user.Preference.Medium,
		PreferenceHighChannel:   user.Preference.High,
	}, nil
}

func DeleteUserFromDB(id string) error {
	//Delete the user from the database
	queries := db.New(DB.Pool)
	// Parse string to UUID first
	pgTypeUUID, err := parseStringIdToPgTypeUUID(id)
	if err != nil {
		return fmt.Errorf("failed to parse UUID: %v", err)
	}

	// Execute delete operation with proper UUID type
	err = queries.DeleteUserByID(context.Background(), pgTypeUUID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}

func parseStringIdToPgTypeUUID(id string) (pgtype.UUID, error) {
	// Parse string to UUID first
	parseUUID, err := uuid.Parse(id)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid UUID format: %v", err)
	}

	// Convert to pgtype.UUID
	pgID := pgtype.UUID{Bytes: parseUUID, Valid: true}

	return pgID, nil
}

func GetUserFromDB(id string) (model.User, error) {
	//Get the user from the database
	queries := db.New(DB.Pool)
	//Parse the string to UUID
	pgTypeUUID, err := parseStringIdToPgTypeUUID(id)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to parse UUID: %v", err)
	}

	//Execute the query
	user, err := queries.GetUserByID(context.Background(), pgTypeUUID)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user: %v", err)
	}

	//Convert the user to model.User
	userModel, err := parseDBUserToUser(user)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to parse user: %v", err)
	}

	return userModel, nil

}

func UUIDToString(t pgtype.UUID) string {
	if !t.Valid {
		return ""
	}
	src := t.Bytes
	return fmt.Sprintf("%x-%x-%x-%x-%x", src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}

func parseDBUserToUser(dbUser db.User) (model.User, error) {
	// Convert the dbUser to model.User
	user := model.User{
		Id:       uuid.MustParse(UUIDToString(dbUser.ID)),
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Phone:    dbUser.Phone.String,
		Status:   dbUser.Status,
		Preference: model.Preference{
			Low:    dbUser.PreferenceLowChannel,
			Medium: dbUser.PreferenceMediumChannel,
			High:   dbUser.PreferenceHighChannel,
		},
	}

	return user, nil
}

func GetAllUsersFromDb() ([]model.User, error) {
	//Get all the users from the database
	queries := db.New(DB.Pool)
	users, err := queries.GetAllUsers(context.Background())
	if err != nil {
		fmt.Println("Failed to get users: ", err)
	}

	var allUsers []model.User
	for _, user := range users {
		userModel, err := parseDBUserToUser(user)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user: %v", err)
		}
		allUsers = append(allUsers, userModel)
	}

	return allUsers, nil

}

package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type User struct {
	ID      int
	Name    string
	Age     int
	Address Address
}

var (
	savePath = "./saved_users"
	nextID   = GetLastUserId(savePath) + 1
)

func AddUser(u User) (User, error) {
	if u.ID != 0 {
		return User{}, errors.New("new user must not include id or it must be set to zero")
	}
	u.ID = nextID
	nextID++
	if err := SaveUserData(savePath, u); err != nil {
		return User{}, err
	}
	return u, nil
}

func GetUserByID(id int) (User, error) {
	var user User

	user, err := GetSavedUser(savePath, id)
	if err != nil {
		return user, fmt.Errorf("user with ID '%v' not found", id)
	}

	return user, nil
}

// New Methods

func GetLastUserId(path string) int {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return len(dirEntries)
}

func SaveUserData(savePath string, user User) error {
	// 1. encode user data as json
	jsonData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshalling user data: %w", err)
	}

	// 2. ensure the save directory exists
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return fmt.Errorf("failed to create save directory %s: %w", savePath, err)
	}

	// 3. create a file with the user_{id}.json then save the json data to it
	fileSavePath := fmt.Sprintf("%s/user_%d.json", savePath, user.ID)

	// go uses standard unix/linux permissions
	// 0644 means (owner read/write, group read, others read)
	err = os.WriteFile(fileSavePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to saved user file %s: %w", fileSavePath, err)
	}

	fmt.Printf("Successfully saved user data to %s\n", fileSavePath)
	return nil
}

func GetSavedUser(savePath string, userId int) (User, error) {
	var user User

	// 1. if the user is saved it should be at user_{id}.json
	fileSavePath := fmt.Sprintf("%s/user_%d.json", savePath, userId)

	// 2. read the file content, if the file doesn't exists return nil
	userJson, err := os.ReadFile(fileSavePath)
	if err != nil {
		// Return the zero value and the wrapped error
		return user, fmt.Errorf("failed to load user file %s: %w", fileSavePath, err)
	}

	err = json.Unmarshal(userJson, &user)
	if err != nil {
		return user, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return user, nil
}

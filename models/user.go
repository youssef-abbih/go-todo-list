package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       	uint   			`json:"id" gorm:"primaryKey"`
	Email    	string 			`gorm:"unique" json:"email"`
	Password 	string 			`json:"password"`
	CreatedAt 	time.Time 		`json:"created_at"`
	UpdatedAt   time.Time      	`json:"updated_at"`
	DeletedAt   gorm.DeletedAt 	`gorm:"index" json:"-"`
}

func GetUsers() []User {
	var users []User
	DB.Find(&users)
	return users
}

func AddUser(user User) (User, error){
	user.CreatedAt = time.Now()
	bytePassword, err :=bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil{
		return User{}, err
	}
	user.Password = string(bytePassword)
	if err := DB.Create(&user).Error; err != nil {
    return User{}, err
}

	return user, nil
}

func GetUserByEmail(email string) (User, bool) {

	var user User

	result := DB.Where("email = ?", email).First(&user)
	
	if result.Error != nil {
		return User{}, false
	}
	return user, true
}

func UpdateUser(id uint, updatedUser User) (User, bool){

	var existingUser User

	result := DB.First(&existingUser, id)

	if result.Error != nil {
		return User{}, false
	}

	if 	existingUser.ID == updatedUser.ID &&
		existingUser.Email == updatedUser.Email &&
		existingUser.Password == updatedUser.Password{
			return existingUser, true
		}
	updatedUser.ID = existingUser.ID
	updatedUser.CreatedAt = existingUser.CreatedAt
	DB.Save(&updatedUser)
	return updatedUser, true
}

func DeleteUser(id uint) (User, bool){
	
	var user User
	result := DB.Unscoped().First(&user, id)

	if result.Error != nil {
		return User{}, false
	}
	if user.DeletedAt.Valid {
		return User{}, false // Already deleted
	}
	
	tasks := GetTasks(user.ID)
	for _, task := range tasks{
		DeleteTask(task.ID, user.ID)
	}
	
	DB.Delete(&user)
	return user, true
}
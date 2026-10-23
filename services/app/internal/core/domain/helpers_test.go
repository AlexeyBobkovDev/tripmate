package domain

func compareExistedUserWithGotUser(existedUser, gotUser *User) bool {
	areEqual := existedUser.ID == gotUser.ID &&
		existedUser.Version == gotUser.Version &&
		existedUser.Name == gotUser.Name &&
		existedUser.Surname == gotUser.Surname &&
		existedUser.Username == gotUser.Username &&
		existedUser.BirthDate.Equal(gotUser.BirthDate) &&
		existedUser.Description == gotUser.Description &&
		existedUser.Email == gotUser.Email &&
		existedUser.PhoneNumber == gotUser.PhoneNumber &&
		existedUser.CreatedAt.Equal(gotUser.CreatedAt) &&
		existedUser.UpdatedAt.Equal(gotUser.UpdatedAt)
	if existedUser.DeletedAt == nil && gotUser.DeletedAt == nil {
		return areEqual &&
			existedUser.DeletedAt == nil &&
			gotUser.DeletedAt == nil
	}
	return areEqual && existedUser.DeletedAt.Equal(*gotUser.DeletedAt)
}

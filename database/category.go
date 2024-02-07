package database

import "gorm.io/gorm"

type CategoryModel struct {
	gorm.Model

	Name        string `gorm:"not null,uniqueIndex"`
	Description string `gorm:"not null,default:''"`
}

/*
* GetCategories returns all categories in the database.
* It returns an error if the categories could not be fetched.
 */
func GetCategories() ([]CategoryModel, error) {
	var categories []CategoryModel

	err := Database.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

/*
* GetCategoryByName returns the category with the specified name.
* It returns an error if the category could not be fetched.
 */
func GetCategoryByName(name string) (CategoryModel, error) {
	var category CategoryModel

	err := Database.Where("name = ?", name).First(&category).Error
	if err != nil {
		return category, err
	}

	return category, nil
}

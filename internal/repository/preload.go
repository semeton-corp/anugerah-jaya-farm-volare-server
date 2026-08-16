package repository

import "gorm.io/gorm"

func preloadWithDeleted(db *gorm.DB) *gorm.DB {
	return db.Unscoped()
}

func preloadCageWithDeleted(query *gorm.DB, path string) *gorm.DB {
	return query.
		Preload(path, preloadWithDeleted).
		Preload(path + ".Location")
}

func preloadChickenCageWithDeleted(query *gorm.DB, path string) *gorm.DB {
	return query.
		Preload(path, preloadWithDeleted).
		Preload(path+".Cage", preloadWithDeleted).
		Preload(path + ".Cage.Location")
}

package main

import "testing"

func TestAddToCollection(t *testing.T) {
	db := InitDB(":memory:")
	CreateTables(db)
	AddToCollection(db, 9999, "TestDrive", 8.8, "Test")
}

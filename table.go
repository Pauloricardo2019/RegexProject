package main

type Table struct {
	ID   uint `gorm:primaryKey`
	IP   string
	Date string
	Verb string
}

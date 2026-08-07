package models

import "errors"

const MaxNameLength int = 140
const MaxDescriptionLength int = 1000

type Order string

const OrderASC Order = "ASC"
const OrderDESC Order = "DESC"

type WeightUnit string

const WeightUnit_Kg WeightUnit = "kg"
const WeightUnit_Lbs WeightUnit = "lbs"
const WeightUnit_Mixed WeightUnit = "mixed"

var InvalidWeightUnit error = errors.New("Invalid weight unit")

const (
	StatusRequested  = "requested"
	StatusConfirmed  = "confirmed"
	StatusRestricted = "restricted"
	StatusBanned     = "banned"
	StatusAdmin      = "admin"
)

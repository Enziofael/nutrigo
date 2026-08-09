package models

import "errors"

// ORDER

type Order string

const (
	Order_ASC  Order = "ASC"
	Order_DESC Order = "DESC"
)

func (o Order) Validate() bool {
	switch o {
	case Order_ASC, Order_DESC:
		return true
	default:
		return false
	}
}

var ErrInvalidOrder error = errors.New("invalid sort order")

// UNIT PREFERENCE

// Valid values [1;23]
type UnitPreferences int32

// SYNC NEW UNITS WITH:
//	- [unitPreferences_max]
//	- [Exercise]
//	- [ExerciseCreateRequest]
//	- [ExercisePatchRequest]
const (
	UnitPreference_Kg    UnitPreferences = 0b00_0_01
	UnitPreference_Lbs   UnitPreferences = 0b00_0_10
	UnitPreference_Time  UnitPreferences = 0b00_1_00
	UnitPreference_Km    UnitPreferences = 0b01_0_00
	UnitPreference_Miles UnitPreferences = 0b10_0_00
)

const (
	unitPreferences_restriction_km_miles = 0b11_0_00
	unitPreferences_min UnitPreferences = UnitPreference_Kg
	unitPreferences_max UnitPreferences = UnitPreference_Kg | UnitPreference_Lbs | UnitPreference_Time | UnitPreference_Km | UnitPreference_Miles
)

func (p UnitPreferences) Validate() bool {
	if p >= unitPreferences_min && p <= unitPreferences_max && !p.Includes(unitPreferences_restriction_km_miles) {
		return true
	}
	return false
}

// Checks if p includes specified unit u
//
// If you're checking for 1 unit and not grouped unit preference
// use shortcuts instead:
//
//	func (p UnitPreference) IncludesKg() bool
//	func (p UnitPreference) IncludesLbs() bool
//	func (p UnitPreference) IncludesTime() bool
//	func (p UnitPreference) IncludesKm() bool
//	func (p UnitPreference) IncludesMile() bool
func (p UnitPreferences) Includes(u UnitPreferences) bool {
	return p&u == u
}

// Shortcut for p.Includes(UnitPreference_Kg)
func (p UnitPreferences) IncludesKg() bool {
	return p.Includes(UnitPreference_Kg)
}

// Shortcut for p.Includes(UnitPreference_Lbs)
func (p UnitPreferences) IncludesLbs() bool {
	return p.Includes(UnitPreference_Lbs)
}

// Shortcut for p.Includes(UnitPreference_Time)
func (p UnitPreferences) IncludesTime() bool {
	return p.Includes(UnitPreference_Time)
}

// Shortcut for p.Includes(UnitPreference_Km)
func (p UnitPreferences) IncludesKm() bool {
	return p.Includes(UnitPreference_Km)
}

// Shortcut for p.Includes(UnitPreference_Mile)
func (p UnitPreferences) IncludesMiles() bool {
	return p.Includes(UnitPreference_Miles)
}

var ErrInvalidUnitPreference error = errors.New("invalid unit preference")

// USERSTATUS

type UserStatus string

const (
	UserStatus_Unknown    = "unknown"
	UserStatus_Requested  = "requested"
	UserStatus_Confirmed  = "confirmed"
	UserStatus_Restricted = "restricted"
	UserStatus_Banned     = "banned"
	UserStatus_Admin      = "admin"
)

func (s UserStatus) Validate() bool {
	switch s {
	case
		UserStatus_Unknown, UserStatus_Requested, UserStatus_Confirmed,
		UserStatus_Restricted, UserStatus_Banned, UserStatus_Admin:
		return true
	default:
		return false
	}
}

var ErrInvalidUserStatus error = errors.New("invalid user status")

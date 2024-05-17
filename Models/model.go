package Models

import (
	"time"
)

type User struct {
	Email        string    `bson:"email,omitempty" json:"email,omitempty"`
	Firstname    string    `bson:"firstname,omitempty" json:"firstname,omitempty"`
	Lastname     string    `bson:"lastname,omitempty" json:"lastname,omitempty"`
	JMBG         string    `bson:"jmbg,omitempty" json:"jmbg,omitempty"`
	Birthdate    string    `bson:"birthdate,omitempty" json:"birthdate,omitempty"`
	Gender       string    `bson:"gender,omitempty" json:"gender,omitempty"`
	Role         string    `bson:"role,omitempty" json:"role,omitempty"`
	Street       string    `bson:"street,omitempty" json:"street,omitempty"`
	StreetNumber string    `bson:"streetnumber,omitempty" json:"streetnumber,omitempty"`
	City         string    `bson:"city,omitempty" json:"city,omitempty"`
	Country      string    `bson:"country,omitempty" json:"country,omitempty"`
	Requests     []Request `bson:"requests,omitempty" json:"requests,omitempty"`
	Weapons      []Weapon  `bson:"weapons,omitempty" json:"weapons,omitempty"`
}

type Weapon struct {
	WeaponType   string `bson:"weapontype,omitempty" json:"weapontype,omitempty"`
	SerialNumber int    `bson:"serialnumber,omitempty" json:"serialnumber,omitempty"`
	Caliber      string `bson:"caliber,omitempty" json:"caliber,omitempty"`
}

type Type int

const (
	Handgun Type = iota
	Rifle
	Carbine
	Axe
	LongBarreledPistol
	PumpActionShotgun
	BowAndArrow
	SelfDefenseWeaponry
)

var typeNames = [...]string{
	"HANDGUN",
	"RIFLE",
	"CARBINE",
	"AXE",
	"LONG-BARRELED PISTOL",
	"PUMP-ACTION SHOTGUN",
	"BOW AND ARROW",
	"SELF-DEFENSE WEAPONRY",
}

type Responsepros struct {
	Prosecuted bool `json:"prosecuted"`
}
type Response struct {
	VehicleWanted bool `json:"vehicle_wanted"`
}

func (wt Type) String() string {
	names := [...]string{"HANDGUN", "RIFLE", "CARBINE", "AXE", "LONG-BARRELED PISTOL", "PUMP-ACTION SHOTGUN", "BOW AND ARROW", "SELF-DEFENSE WEAPONRY"}
	if wt < Handgun || wt > SelfDefenseWeaponry {
		return "Unknown"
	}
	return names[wt]
}

type Date time.Time

func (d Date) String() string {
	return time.Time(d).Format("2006-01-02") // Formatirajte datum kako želite
}

// SpecialCross or temporaryExport
type Request struct {
	Uuid         string `bson:"uuid,omitempty" json:"uuid,omitempty"`
	RequestDate  string `bson:"requestdate,omitempty" json:"requestdate,omitempty"`
	RequestState string `bson:"requeststate,omitempty" json:"requeststate,omitempty"`
	Weapon       Weapon `bson:"weapon,omitempty" json:"weapon,omitempty"`
	Email        string `bson:"email,omitempty" json:"email,omitempty"`
	Recorded     string `bson:"recorded,omitempty" json:"recorded,omitempty"`
}
type Vehicle struct {
	Plates string `bson:"plates,omitempty" json:"plates,omitempty"`
}

type GetRequest struct {
	Uuid string
}

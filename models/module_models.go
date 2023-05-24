package models

type Module struct {
	ModuleName string `json:"module-name,omitempty" validate:"required"`
	LevelList  []Level
}

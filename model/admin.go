package model
type Admin struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name		 string `gorm:"not null" json:"name"`
	PasswordHash string `gorm:"not null" json:"password_hash,omitempty"`

}
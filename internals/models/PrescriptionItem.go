package models

// PrescriptionItem links prescriptions to medicines
// @Description Each prescription item specifies medicine and quantity
type PrescriptionItem struct {
	ID             uint `json:"id"`
	PrescriptionID uint `json:"prescription_id"` // FK -> Prescription.ID
	MedicineID     uint `json:"medicine_id"`     // FK -> Medicine.ID
	Quantity       int  `json:"quantity"`
}
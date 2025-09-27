package models

// Prescription represents a prescription created by a doctor
// @Description Contains patient info and prescribed medicine details
type Prescription struct {
	ID uint `json:"id"`
	PatientName string `json:"patientname"`
	MedicineName string `json:"medicinename"`
	Dosage string `json:"dosage"`
	Quantity int `json:"quantity"`
}

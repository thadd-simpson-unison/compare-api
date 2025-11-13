package Models

type CompareResponse struct {
	SuppliedEmailsCount int32
	RedEmailCount       int32
	BlueEmailCount      int32
	RedRequestsDuration float64
	BlueRequestDuration float64
	RedEmailFails       []string
	BlueEmailFails      []string
	CombinedResults     []CombinedResult
}

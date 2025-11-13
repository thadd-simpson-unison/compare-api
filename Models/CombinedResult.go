package Models

type CombinedResult struct {
	Source      string
	Email       string
	Status      string
	RedResponse *RedResponse
	BlueData    *[]string
}

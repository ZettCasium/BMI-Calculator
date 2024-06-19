package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// BMIRequest untuk menangkap request dari client
type BMIRequest struct {
	Age           int     `json:"age"`
	Weight        float64 `json:"weight"`        // dalam kg
	Height        float64 `json:"height"`        // dalam cm
	ActivityLevel float64 `json:"activityLevel"` // Tambahkan ini
}

// BMIResponse untuk mengirimkan response ke client
type BMIResponse struct {
	BMI           float64 `json:"bmi"`
	Category      string  `json:"category"`
	ActivityLevel float64 `json:"caloricNeeds"` // Tambahkan ini
}

// hitungBMI menghitung BMI dan kategorinya
func calculateBMI(w http.ResponseWriter, r *http.Request) {
	var req BMIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received data: %+v", req)

	log.Printf("Received activityLevel: %f", req.ActivityLevel)

	bmi := req.Weight / ((req.Height / 100) * (req.Height / 100))
	category := ""
	switch {
	case bmi < 18.5:
		category = "Underweight"
	case bmi < 24.9:
		category = "Normal weight"
	case bmi < 29.9:
		category = "Overweight"
	default:
		category = "Obesity"
	}

	ageFloat := float64(req.Age)
	bmr := 10*req.Weight + 6.25*req.Height - 5*ageFloat + 5

	log.Printf("Calculated BMR: %f", bmr)

	caloricNeeds := bmr * req.ActivityLevel

	log.Printf("Calculated Caloric Needs: %f", caloricNeeds)

	resp := BMIResponse{
		BMI:           bmi,
		Category:      category,
		ActivityLevel: caloricNeeds, // Kirimkan kembali tingkat aktivitas
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

/*func hitungBMI(weight, height float64) (float64, string) {
	heightInMeters := height / 100
	bmi := weight / math.Pow(heightInMeters, 2)

	var category string
	switch {
	case bmi < 18.5:
		category = "Underweight"
	case bmi >= 18.5 && bmi < 24.9:
		category = "Normal weight"
	case bmi >= 25 && bmi < 29.9:
		category = "Overweight"
	default:
		category = "Obese"
	}

	return bmi, category
}*/

// Middleware CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Jika OPTIONS (preflight), balas dengan 200 dan hentikan.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Lanjutkan ke handler berikutnya.
		next.ServeHTTP(w, r)
	})
}

/*func bmiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request BMIRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bmi, category := hitungBMI(request.Weight, request.Height)
	response := BMIResponse{
		BMI:      bmi,
		Category: category,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}*/

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/calculate-bmi", calculateBMI)

	// Tambahkan middleware CORS
	corsHandler := corsMiddleware(mux)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

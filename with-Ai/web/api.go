package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"gorgonia.org/tensor"
)

type Request struct {
	X float64 `json:"X"`
	P float64 `json:"P"`
	S float64 `json:"S"`
	V float64 `json:"V"`
}

type Response struct {
	PH float64 `json:"pH"`
}

// Load trained LSTM model
var lstm *LSTM

func predictHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert input to tensor
	inputTensor := tensor.New(tensor.WithShape(1, 4), tensor.Of(tensor.Float32),
		tensor.WithBacking([]float32{float32(req.X), float32(req.P), float32(req.S), float32(req.V)}))

	// Predict pH
	predictedPH := lstm.Predict(inputTensor)

	resp := Response{PH: predictedPH}
	json.NewEncoder(w).Encode(resp)
}

func main() {
	lstm = NewLSTM()
	lstm.LoadModel()

	http.HandleFunc("/predict", predictHandler)

	fmt.Println("🚀 Go API running at: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

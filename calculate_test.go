package main

import (
	"math"
	"testing"
)

func TestCalculateMonthlyPayment(t *testing.T) {
	type args struct {
		amount int64
		term   int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "Calculate monthly payment for $10,000 loan over 3 years",
			args: args{amount: 10000, term: 36},
			want: 299,
		},
		{
			name: "Calculate monthly payment for $25,000 loan over 5 years",
			args: args{amount: 25000, term: 60},
			want: 471,
		},
		{
			name: "Calculate monthly payment for $50,000 loan over 10 years",
			args: args{amount: 50000, term: 120},
			want: 530,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateMonthlyPayment(tt.args.amount, tt.args.term); got != tt.want {
				t.Errorf("calculateMonthlyPayment() = %v, want %v", got, tt.want)
			}
		})
	}

	// Test for edge case: large loan amount and term
	largeAmount := int64(math.Pow(10, 9)) // 1 billion dollars
	largeTerm := int64(math.Pow(10, 4))   // 10,000 months
	want := int64(4166666)                // expected monthly payment
	got := calculateMonthlyPayment(largeAmount, largeTerm)
	if got != want {
		t.Errorf("calculateMonthlyPayment() = %v, want %v", got, want)
	}
}

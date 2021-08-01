package datatypes

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLocation_MarshalText(t *testing.T) {
	type Sample struct {
		Location Location
		Result   string
	}

	height := float64(100)
	sampleList := []Sample{
		{Location: Location{Point: [2]float64{101.1, 32.1}}, Result: "101.1,32.1"},
		{Location: Location{Point: [2]float64{101.1, 32.1}, height: &height}, Result: "101.1,32.1,100"},
	}

	for i := 0; i < len(sampleList); i++ {
		res, err := sampleList[i].Location.MarshalText()
		require.NoError(t, err)
		require.Equal(t, sampleList[i].Result, string(res))
	}
}

func TestLocation_UnmarshalText(t *testing.T) {
	type Sample struct {
		Location string
		Result   Location
	}

	height := float64(100)
	sampleList := []Sample{
		{Location: "101.1,32.1", Result: Location{Point: [2]float64{101.1, 32.1}}},
		{Location: "101.1,32.1,100", Result: Location{Point: [2]float64{101.1, 32.1}, height: &height}},
	}

	for i := 0; i < len(sampleList); i++ {
		location := Location{}
		err := location.UnmarshalText([]byte(sampleList[i].Location))
		require.NoError(t, err)
		require.Equal(t, sampleList[i].Result, location)
	}
}

package datatypes

import (
	"errors"
	"github.com/go-courier/geography"
	"strconv"
	"strings"
)

// openapi:strfmt location
type Location struct {
	geography.Point
	height *float64
}

func (location Location) MarshalText() ([]byte, error) {
	return []byte(location.String()), nil
}

func (location *Location) UnmarshalText(data []byte) (err error) {
	l := strings.Split(string(data), ",")
	if len(l) != 2 && len(l) != 3 {
		return errors.New("Location UnmarshalText error, param error")
	}
	if len(l) == 3 {
		height, err := strconv.ParseFloat(l[2], 64)
		if err != nil {
			return err
		}
		location.height = &height
	}

	if location.Point[0], err = strconv.ParseFloat(l[0], 64); err != nil {
		return err
	}
	if location.Point[1], err = strconv.ParseFloat(l[1], 64); err != nil {
		return err
	}

	return nil
}

func (location Location) String() string {
	res := ""
	res = strconv.FormatFloat(location.Point[0], 'f', -1, 64)
	res = res + "," + strconv.FormatFloat(location.Point[1], 'f', -1, 64)
	if location.height != nil {
		res = res + "," + strconv.FormatFloat(*location.height, 'f', -1, 64)
	}
	return res
}

type Locations []Location

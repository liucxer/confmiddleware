package datatypes

import (
	"fmt"
	"time"
	_ "time/tzdata"

	"github.com/go-courier/sqlx/v2/builder"
	sqlxDatatypes "github.com/go-courier/sqlx/v2/datatypes"
)

// openapi:strfmt date-time-or-range
type DateTimeOrRange struct {
	From sqlxDatatypes.Timestamp
	To   sqlxDatatypes.Timestamp
	ValueOrRangeOpt
}

func (timeRange *DateTimeOrRange) IsZero() bool {
	return timeRange == nil || (timeRange.From.IsZero() && timeRange.To.IsZero())
}

func (timeRange *DateTimeOrRange) ConditionFor(c *builder.Column) builder.SqlCondition {
	return timeRange.ValueOrRangeOpt.ConditionFor(c, &timeRange.From, &timeRange.To)
}

func (timeRange DateTimeOrRange) MarshalText() ([]byte, error) {
	return timeRange.ValueOrRangeOpt.MarshalText(&timeRange.From, &timeRange.To)
}

func (timeRange *DateTimeOrRange) UnmarshalText(data []byte) error {
	tr := DateTimeOrRange{}

	err := tr.ValueOrRangeOpt.UnmarshalText(data, &tr.From, &tr.To)
	if err != nil {
		return err
	}

	if !tr.From.IsZero() && !tr.To.IsZero() {
		if time.Time(tr.From).After(time.Time(tr.To)) {
			return fmt.Errorf("time from should not after time to")
		}
	}

	*timeRange = tr
	return nil
}

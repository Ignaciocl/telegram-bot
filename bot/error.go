package bot

import "errors"

var (
	ErrNoUserScheduled        = errors.New("no user scheduled news for the given hour")
	errRetrievingScheduleInfo = errors.New("error retrieving schedule information")
)

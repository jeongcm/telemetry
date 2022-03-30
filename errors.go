package telemetry

import "errors"

var (
	// ErrNotFoundService 는 telemetry service를 조회했으나 찾지 못했을 때 발생한다.
	ErrNotFoundService = errors.New("not found telemetry service")
	// ErrConnectionFailed 는 telemetry service에 연결될 새로운 Client가 만들어지지 못했을때 발생한다.
	ErrConnectionFailed = errors.New("could not connect to telemetry service")
)

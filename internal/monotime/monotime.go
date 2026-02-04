package monotime
import(
	"time"
)

type TimeSource interface{
	Now() time.Time
}

type RealTimeSource struct{}

func(t *RealTimeSource) Now() time.Time{
	return time.Now().UTC()
}
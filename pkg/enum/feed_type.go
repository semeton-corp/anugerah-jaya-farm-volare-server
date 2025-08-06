package enum

type FeedType uint8

const (
	FeedTypeUnknown        FeedType = 0
	FeedTypeReadyToEatFeed FeedType = 1
	FeedTypeRawFeed        FeedType = 2
)

var (
	FeedTypeMap = map[FeedType]string{
		FeedTypeReadyToEatFeed: "Pakan Jadi",
		FeedTypeRawFeed:        "Pakan Adukan",
	}
)

func (c FeedType) String() string {
	return FeedTypeMap[c]
}

func ValueOfFeedType(value string) FeedType {
	for k, v := range FeedTypeMap {
		if v == value {
			return k
		}
	}
	return FeedTypeUnknown
}

func (c FeedType) IsValid() bool {
	_, ok := FeedTypeMap[c]
	return ok
}

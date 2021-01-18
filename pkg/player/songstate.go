package player

type SongState int

const (
	Stopped SongState = iota
	Playing
	Paused
)

func (s SongState) String() string {
	return [...]string{"Stopped", "Playing", "Paused"}[s]
}

package tui

type status int

const (
	LoadingView status = iota
	VaultListView
	LoginsListView
	LoginDetailView
)

type Board struct {
	focus status
}

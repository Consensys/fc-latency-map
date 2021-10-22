package webhook

type Notifier interface {
	Notify(files *[]string) bool
}

package sync

type Callback interface {
	OnSuccess()
	OnFailure(error)
}

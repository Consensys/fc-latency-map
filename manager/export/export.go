package export

//go:generate mockgen -destination mocks.go -package export . Service

type Service interface {

	// export from db to json file
	export()
}

package export

//go:generate mockgen -destination mocks.go -package export . Service

type Service interface {

	// dbExportData from db to json file
	export(fn string)
}

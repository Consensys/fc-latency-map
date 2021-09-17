package export

type Service interface {

	// dbExportData from db to json file
	export(fn string)
}

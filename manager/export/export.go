package export

type ExportService interface {

	// dbExportData from db to json file
	export(fn string)
}

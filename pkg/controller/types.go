package controller

type Controller struct {
	Timeout  int
	Parallel int
	Pulls    int
}

type Report struct {
	Image            string
	Registry         string
	FailedLayers     int
	DownloadedLayers int
	AvgDownloadTime  float64
	Success          bool
	ManifestDownload bool
}

type ReportsList struct {
	Reports []Report
}

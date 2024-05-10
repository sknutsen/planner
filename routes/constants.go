package routes

const (
	Index     string = "/"
	IndexWeek string = "/:week"
	Assets    string = "/assets"
	Day       string = "day/:date"
	Daytasks  string = Day + "/tasks"
)

package routes

const (
	Index     string = "/"
	IndexWeek string = "/:week"

	Assets string = "/assets"

	Day      string = "/day/:date"
	Daytasks string = Day + "/tasks"

	Components      string = "/components"
	ComponentsModal string = Components + "/modal"

	Task       string = "/task/:id"
	TaskCreate string = Day + "/create"
	TaskDelete string = Task + "/delete"
	TaskEdit   string = Task + "/edit"
	TaskUpdate string = "/task/update"
)

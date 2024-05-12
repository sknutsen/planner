package routes

const (
	Index     string = "/"
	IndexWeek string = "/:week"

	Assets string = "/assets"

	Callback string = "/callback"

	Components      string = "/components"
	ComponentsModal string = Components + "/modal"

	Day      string = "/day/:date"
	Daytasks string = Day + "/tasks"

	Login  string = "/login"
	Logout string = "/logout"

	Plan       string = "/plan/:id"
	PlanCreate string = "/plan/create"
	PlanDelete string = Plan + "/delete"
	PlanEdit   string = Plan + "/edit"
	PlanUpdate string = "/plan/update"

	Task       string = "/task/:id"
	TaskCreate string = Day + "/create"
	TaskDelete string = Task + "/delete"
	TaskEdit   string = Task + "/edit"
	TaskUpdate string = "/task/update"
)

package routes

const (
	Index string = "/"

	Assets string = "/assets"

	Callback string = "/callback"

	Components      string = "/components"
	ComponentsModal string = Components + "/modal"

	Day      string = "/:planId/day/:date"
	Daytasks string = Day + "/tasks"

	History          string = "/history"
	HistoryPlan      string = History + "/:planId"
	HistoryPlanTasks string = History + "/:planId/tasks"

	Login  string = "/login"
	Logout string = "/logout"

	Plan       string = "/plan/:id"
	PlanCreate string = "/plan/create"
	PlanDelete string = Plan + "/delete"
	PlanEdit   string = Plan + "/edit"
	PlanUpdate string = "/plan/update"

	Resources              string = "/resources"
	ResourcesPlan          string = Resources + "/:planId"
	ResourcesPlanResources string = Resources + "/:planId/resources"
	Resource               string = "/resource/:id"
	ResourceCreate         string = ResourcesPlan + "/create"
	ResourceDelete         string = Resource + "/delete"
	ResourceEdit           string = Resource + "/edit"
	ResourceUpdate         string = "/resource/update"

	Task         string = "/task/:id"
	TaskComplete string = Task + "/complete"
	TaskCopy     string = "/task/copy"
	TaskCreate   string = Day + "/create"
	TaskDelete   string = Task + "/delete"
	TaskEdit     string = Task + "/edit"
	TaskUpdate   string = "/task/update"

	User string = "/user"

	Week         string = "/week"
	WeekPlan     string = Week + "/:id"
	WeekPlanWeek string = WeekPlan + "/:week"
)

package routes

const (
	// args
	date   string = "/:date"
	id     string = "/:id"
	planId string = "/:planId"
	weekNo string = "/:week"

	// actions
	complete string = "/complete"
	copy     string = "/copy"
	create   string = "/create"
	delete   string = "/delete"
	edit     string = "/edit"
	update   string = "/update"

	// paths
	Index string = "/"

	Assets string = "/assets"

	Callback string = "/callback"

	Components      string = "/components"
	modal           string = "/modal"
	ComponentsModal string = Components + modal

	day      string = "/day"
	Day      string = planId + day + date
	Daytasks string = Day + tasks

	History          string = "/history"
	HistoryPlan      string = History + planId
	HistoryPlanTasks string = History + planId + tasks

	Login  string = "/login"
	Logout string = "/logout"

	plan       string = "/plan"
	Plan       string = plan + id
	PlanCreate string = plan + create
	PlanDelete string = Plan + delete
	PlanEdit   string = Plan + edit
	PlanUpdate string = plan + update

	Resources              string = "/resources"
	ResourcesPlan          string = Resources + planId
	ResourcesPlanResources string = Resources + planId + Resources

	resource       string = "/resource"
	Resource       string = resource + id
	ResourceCreate string = ResourcesPlan + create
	ResourceDelete string = Resource + delete
	ResourceEdit   string = Resource + edit
	ResourceUpdate string = resource + update

	tasks              string = "/tasks"
	task               string = "/task"
	Task               string = task + id
	TaskComplete       string = Task + complete
	TaskCopy           string = task + copy
	TaskCreate         string = Day + create
	TaskCreateTemplate string = TaskCreate + template
	TaskDelete         string = Task + delete
	TaskEdit           string = Task + edit
	TaskUpdate         string = task + update

	Templates              string = "/templates"
	TemplatesPlan          string = Templates + planId
	TemplatesPlanTemplates string = Templates + planId + Templates

	template       string = "/template"
	Template       string = template + id
	TemplateCreate string = TemplatesPlan + create
	TemplateDelete string = Template + delete
	TemplateEdit   string = Template + edit
	TemplateUpdate string = template + update

	User string = "/user"

	Week         string = "/week"
	WeekPlan     string = Week + id
	WeekPlanWeek string = WeekPlan + weekNo
)

package exports

// Comment represents the exported version of a GitHub comment
type Comment struct {
	ID               int    `json:"id"`
	Body             string `json:"body"`
	Position         int    `json:"position"`
	OriginalPosition int    `json:"original_position"`
	User             *User  `json:"user"`
}

// Pull represents the exported version of a Github pull request
type Pull struct {
	ID     int    `json:"id"`
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Merged bool   `json:"merged"`
	User   *User  `json:"user"`
	Repo   *Repo  `json:"repo"`
}

// User represents the exported version of a user in GitHub
type User struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

// Repo represents the exported version of a GitHub repository
type Repo struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Description   string `json:"description"`
	WatchersCount int    `json:"watchers_count"`
	Language      string `json:"language"`
	Owner         *Owner `json:"owner"`
}

// Owner represents the exported version of a GitHub repository
type Owner struct {
	ID int `json:"id"`
}

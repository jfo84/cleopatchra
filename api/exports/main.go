package exports

// Comment represents the exported version of a GitHub comment
type Comment struct {
	ID               int    `jsonapi:"primary,comments"`
	Body             string `jsonapi:"attr,body"`
	Position         int    `jsonapi:"attr,position"`
	OriginalPosition int    `jsonapi:"attr,original_position"`
	User             *User  `jsonapi:"relation,user"`
}

// Pull represents the exported version of a Github pull request
type Pull struct {
	ID             int        `jsonapi:"primary,pulls"`
	Number         int        `jsonapi:"attr,number"`
	Additions      int        `jsonapi:"attr,additions"`
	Deletions      int        `jsonapi:"attr,deletions"`
	NumComments    int        `jsonapi:"attr,num_comments"`
	Title          string     `jsonapi:"attr,title"`
	Body           string     `jsonapi:"attr,body"`
	MergeableState string     `jsonapi:"attr,mergeable_state"`
	Merged         bool       `jsonapi:"attr,merged"`
	User           *User      `jsonapi:"relation,user"`
	Repo           *Repo      `jsonapi:"relation,repo"`
	Comments       []*Comment `jsonapi:"relation,comments"`
}

// User represents the exported version of a user in GitHub
type User struct {
	ID    int    `jsonapi:"primary,users"`
	Login string `jsonapi:"attr,login"`
}

// Repo represents the exported version of a GitHub repository
type Repo struct {
	ID            int    `jsonapi:"primary,repos"`
	Name          string `jsonapi:"attr,name"`
	FullName      string `jsonapi:"attr,full_name"`
	Description   string `jsonapi:"attr,description"`
	WatchersCount int    `jsonapi:"attr,watchers_count"`
	Language      string `jsonapi:"attr,language"`
	Owner         *Owner `jsonapi:"relation,owner"`
}

// Owner represents the exported version of a GitHub repository
type Owner struct {
	ID int `jsonapi:"primary,owners"`
}

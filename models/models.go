package models

type Threads []Thread

type Error struct {
	Message string `json:"message"`
}

type Author struct {
	Id       int
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	About    string `json:"about"`
}
type Forum struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Author  string `json:"user"`
	Posts   int    `json:"posts"`
	Threads int    `json:"threads"`
}
type Thread struct {
	Id      int    `json:"id"`
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
	Votes   int    `json:"votes"`
}
type Vote struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}

type JsonVote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}

type Post struct {
	Id       int    `json:"id"`
	Author   string `json:"author"`
	Created  string `json:"created"`
	Message  string `json:"message"`
	Forum    string `json:"forum"`
	Thread   int    `json:"thread"`
	IsEdited bool   `json:"isEdited"`
	Parent   int    `json:"parent"`
}

type JsonPost struct {
	Author  string `json:"author"`
	Message string `json:"message"`
	Parent  int    `json:"parent"`
}

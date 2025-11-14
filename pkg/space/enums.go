package space

type SubjectType int

const (
	SubjectUser SubjectType = 1
	
	SubjectMessageBoard SubjectType = iota * 10
	SubjectArticle
	SubjectComment
)

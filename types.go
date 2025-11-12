package main

type JiraResponse struct {
	Data       Data       `json:"data"`
	Extensions Extensions `json:"extensions"`
}

type Data struct {
	Me   Me   `json:"me"`
	Jira Jira `json:"jira"`
}

type Jira struct {
	IssueSearchStable IssueSearchStable `json:"issueSearchStable"`
}

type IssueSearchStable struct {
	TotalCount int    `json:"totalCount"`
	Edges      []Edge `json:"edges"`
}

type Edge struct {
	Node IssueNode `json:"node"`
}

type IssueNode struct {
	IssueId    string     `json:"issueId"`
	WebUrl     string     `json:"webUrl"`
	Fields     Fields     `json:"fields"`
	Key        string     `json:"key"`
	IsResolved bool       `json:"isResolved"`
}

type Fields struct {
	TotalCount int         `json:"totalCount"`
	Edges      []FieldEdge `json:"edges"`
}

type FieldEdge struct {
	Node FieldNode `json:"node"`
}

type FieldNode struct {
	Typename     string  `json:"__typename"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Id           string  `json:"id"`
	FieldId      string  `json:"fieldId"`
	AliasFieldId *string `json:"aliasFieldId"`
	Type         string  `json:"type"`
	User         *User   `json:"user,omitempty"`
	DateTime     *string `json:"dateTime,omitempty"`
}

type Me struct {
	User User `json:"user"`
}

type User struct {
	AccountId     string `json:"accountId"`
	AccountStatus string `json:"accountStatus"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type Extensions struct {
	Gateway Gateway `json:"gateway"`
}

type Gateway struct {
	RequestId       string `json:"request_id"`
	TraceId         string `json:"trace_id"`
	CrossRegion     bool   `json:"crossRegion"`
	EdgeCrossRegion bool   `json:"edgeCrossRegion"`
}

type Variable struct {
	KeyOrId string `json:"keyOrId"`
}

type GraphQlQuery struct {
	Query     string   `json:"query"`
	Variables Variable `json:"variables"`
}

type TextNode struct {
	Attribute *TextNode `json:"attrs"`
	Content   *TextNode `json:"content"`
	Type      string    `json:"type"`
	Text      string    `json:"text"`
	LocalId   string    `json:"localId"`
	State     string    `json:"state"`
}

// StatusInfo represents the current status of an issue
type StatusInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"statusCategory"`
}

// Transition represents an available status transition
type Transition struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	To        StatusInfo `json:"to"`
	IsForward bool       // Derived field for icon display
}

// TransitionsResponse represents the API response for available transitions
type TransitionsResponse struct {
	Transitions []Transition `json:"transitions"`
}

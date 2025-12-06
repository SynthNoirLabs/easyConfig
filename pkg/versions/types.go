package versions

// CommitInfo represents a single commit in the history of a file.
type CommitInfo struct {
	Hash    string `json:"hash"`
	Author  string `json:"author"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

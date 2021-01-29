package ana

// Result is a word that can be a result
type Result string

// Results a list of the result words
type Results []Result

func (a Results) Len() int           { return len(a) }
func (a Results) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Results) Less(i, j int) bool { return len(a[i]) < len(a[j]) }
func (a Results) String() string {
	retTxt := ""
	nl := ""
	for _, v := range a {
		retTxt += string(nl) + string(v)
		nl = "\n"
	}
	return retTxt
}

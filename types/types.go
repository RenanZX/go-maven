package types

type PackageMaven struct {
	GroupId    string
	ArtifactId string
	Version    string
}

func IsEmpty(p PackageMaven) bool {
	return (p.GroupId == "" && p.ArtifactId == "" && p.Version == "")
}

type ResponseServer struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
		Params struct {
			Q               string `json:"q"`
			Core            string `json:"core"`
			DefType         string `json:"defType"`
			Qf              string `json:"qf"`
			Indent          string `json:"indent"`
			Spellcheck      string `json:"spellcheck"`
			Fl              string `json:"fl"`
			Start           string `json:"start"`
			SpellcheckCount string `json:"spellcheck.count"`
			Sort            string `json:"sort"`
			Rows            string `json:"rows"`
			Wt              string `json:"wt"`
			Version         string `json:"version"`
		} `json:"params"`
	} `json:"responseHeader"`
	Response struct {
		NumFound int `json:"numFound"`
		Start    int `json:"start"`
		Docs     []struct {
			ID            string   `json:"id"`
			G             string   `json:"g"`
			A             string   `json:"a"`
			LatestVersion string   `json:"latestVersion"`
			RepositoryID  string   `json:"repositoryId"`
			P             string   `json:"p"`
			Timestamp     int64    `json:"timestamp"`
			VersionCount  int      `json:"versionCount"`
			Text          []string `json:"text"`
			Ec            []string `json:"ec"`
		} `json:"docs"`
	} `json:"response"`
	Spellcheck struct {
		Suggestions []interface{} `json:"suggestions"`
	} `json:"spellcheck"`
}

package entities

type WebPage struct {
	Identifier string
	PrintName  string
	URL        string
}

func NewWebPage(identifier, printName, url string) WebPage {
	return WebPage{
		Identifier: identifier,
		PrintName:  printName,
		URL:        url,
	}
}

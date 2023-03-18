package main

type CurrencyRate struct {
	Fullname   string `xml:"fullname" json:"fullname"`
	Title      string `xml:"title" json:"title"`
	Desciption string `xml:"description" json:"description"`
}

type Rates struct {
	Items []CurrencyRate `xml:"item"`
}

const fakeRssResp = `
<?xml version="1.0" encoding="UTF-8"?>
<rates>
    <item>
		<fullname>Austrailian dollar</fullname>
		<title>AUD</title>
		<description>267.39</description>
    </item>
    <item>
		<fullname>American dollar</fullname>
		<title>USD</title>
		<description>467.45</description>
    </item>
	<item>
		<fullname>EU</fullname>
		<title>EUR</title>
		<description>500.45</description>
    </item>
</rates>`

func rssMock(date string) ([]byte, error) {
	return []byte(fakeRssResp), nil
}

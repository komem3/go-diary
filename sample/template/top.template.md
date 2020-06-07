# diary record{{$base := .Base}}
{{range $i, $yearElem := .Years}}
## {{$yearElem.Year}}
{{range $j, $monthElem := $yearElem.Months}}
## {{$monthElem.Month}}
{{range $k, $dayElem := $monthElem.Days}}
- [{{$dayElem.Day}}]({{$base}}/{{$dayElem.Path}}){{end}}
{{end}}{{end}}

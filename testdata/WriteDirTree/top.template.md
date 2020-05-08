# diary record{{$base := .Base}}
{{range $i, $yearElem := .Years}}
## [{{$yearElem.Year}}]({{$base}}/{{$yearElem.Year}})
{{range $j, $monthElem := $yearElem.Months}}
- [{{$monthElem.Month}}]({{$base}}/{{$yearElem.Year}}/{{$monthElem.Month}})
{{end}}
{{end}}

# diary record
{{range $i, $yearElem := .Years}}
## [{{$yearElem.Year}}](./{{$yearElem.Year}})
{{range $j, $monthElem := $yearElem.Months}}
- [{{$monthElem.Month}}](./{{$yearElem.Year}}/{{$monthElem.Month}})
{{end}}
{{end}}

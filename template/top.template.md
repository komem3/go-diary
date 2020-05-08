# diary record{{$base := .Base}}
{{range $i, $yearElem := .Years}}
## {{$yearElem.Year}}
{{range $j, $monthElem := $yearElem.Months}}
<details>
<summary>{{$monthElem.Month}}</summary>
<ul>{{range $k, $dayElem := $monthElem.Days}}
<li><a href="{{$base}}/{{$dayElem.Path}}">{{$dayElem.Day}}</a></li>{{end}}<ul></details>
{{end}}{{end}}

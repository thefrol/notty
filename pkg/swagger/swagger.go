// Это пакет для генерации сваггер-интерфейса из спеки. Подразумевается, что
// у есть спека формата openapi3 (getkin/kin-openapi). Доступ к такой спеке есть
// у нас поскольку мы пользуемся oapi-codegen
package swagger

import (
	"html/template" //!!!!
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

type SwagData struct {
	SwaggerSpec *openapi3.T
	PageTitle   string
}

func Handler(spec *openapi3.T, title string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		t, err := template.New("swagger").Parse(html)
		if err != nil {
			http.Error(w, "Cant parse html template", http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, SwagData{SwaggerSpec: spec, PageTitle: title})
		if err != nil {
			http.Error(w, "Cant execute template", http.StatusInternalServerError)
			return
		}

	}
}

const html = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<meta
name="description"
content="SwaggerUI"
/>
<title>{{.PageTitle}}</title>
<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-bundle.js" crossorigin></script>
<script>
window.onload = () => {
window.ui = SwaggerUIBundle({
spec: {{.SwaggerSpec}},
dom_id: '#swagger-ui',
});
};
</script>
</body>
</html>
`

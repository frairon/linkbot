package link

import (
	"github.com/frairon/botty"
)

func Home() LinkState {

	return botty.NewStateBuilder[*State]().OnActivate(func(bs LinkSession) {

		type categoryView = struct {
			Title string
			Count int
		}

		catsModel, err := bs.State().st.ListCategories(int64(bs.UserId()))

		var cats []categoryView
		for _, cat := range catsModel {
			cats = append(cats, categoryView{
				Title: cat.Category,
				Count: cat.Count,
			})
		}

		bs.SendError(err)

		bs.SendTemplateMessage(`
Categories:
{{- range $cat := .cats }}
- {{$cat.Title}} ({{$cat.Count}})
{{end }}
		`, botty.TplValues(botty.KV("cats", cats)),
			botty.SendMessageWithKeyboard(botty.NewButtonKeyboard(
				botty.ButtonRow{
					"back",
				},
			),
			))
	}).Build()
}

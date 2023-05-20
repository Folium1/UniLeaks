package delivery

import (
	"github.com/gin-gonic/gin"
)

var foreignLangSubjectsMap = map[string]subjects{
	"eng-gram":  {"Англійська граматика", "eng-gram"},
	"eng-speak": {"Англійська лексика", "eng-speak"},
	"eng-read":  {"Англійське домашнє читання", "eng-read"},
	"eng-phon":  {"Англійська фонетика", "eng-phon"},
	"ger-gram":  {"Німецька граматика", "ger-gram"},
	"ger-speak": {"Німецька лексика", "ger-speak"},
	"ger-read":  {"Німецьке домашнє читання", "ger-read"},
	"izl":       {"Історія зарубіжної літератури", "izl"},
}

func (h Handler) foreignMainPage(c *gin.Context) {
	h.handleSubjectsPage(c, foreignLangSubjectsMap)
}

func (h Handler) foreignSubjects(c *gin.Context) {
	// subject, ok := foreignLangSubjectsMap[c.Param("subj")]
	// if !ok {
	// 	c.AbortWithError(NotFound, errors.New("Not found"))
	// }
	h.tmpl.ExecuteTemplate(c.Writer, "subjectPage.html", foreignLangSubjectsMap)
	// TODO:Leak realization
}

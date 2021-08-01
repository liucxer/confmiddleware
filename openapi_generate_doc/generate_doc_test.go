package openapi_generate_doc_test

import (
	"github.com/liucxer/confmiddleware/openapi_generate_doc"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateDoc(t *testing.T) {
	//url := "http://srv-user-center.base.d.rktl.work/user-center"
	//url := "http://srv-workflow.base.d.rktl.work/workflow"
	//url := "http://srv-bff-iep.intelliep.d.rktl.work/bff-iep"
	//url := "http://srv-database.intelliep.d.rktl.work/database"
	//url := "http://srv-event.intelliep--staging.d.rktl.work/event"
	//url := "http://srv-common-cache.int.querycap.com/common-cache"
	//url := "http://srv-weather.intelliep.d.rktl.work/weather"
	url := "http://srv-iep-permission.intelliep.d.rktl.work/iep-permission"
	openApi, err := openapi_generate_doc.NewOpenApi(url)
	require.NoError(t, err)

	err = openApi.GenerateDoc()
	require.NoError(t, err)

}

//if _, ok := param.Schema.AllOf[0].Reference.Refer.(*oas.ComponentRefer); ok {
//row.AddCell().AddParagraph().AddRun().AddText(param.Schema.AllOf[1].Description)
//} else {
//row.AddCell().AddParagraph().AddRun().AddText(param.Schema.Description)
//}

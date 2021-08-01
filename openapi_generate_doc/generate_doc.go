package openapi_generate_doc

import (
	"encoding/json"
	"fmt"
	"github.com/go-courier/oas"
	"github.com/unidoc/unioffice/color"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/schema/soo/wml"
	"io/ioutil"
	"log"
	"net/http"
)

type OpenApi struct {
	api *oas.OpenAPI
	doc *document.Document
}

func NewOpenApi(url string) (*OpenApi, error) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	openAPI := &oas.OpenAPI{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, openAPI); err != nil {
		return nil, err
	}

	doc := document.New()

	return &OpenApi{
		api: openAPI,
		doc: doc,
	}, nil
}

func (openApi *OpenApi) GenerateDoc() error {
	nd := openApi.doc.Numbering.AddDefinition()
	for url, path := range openApi.api.Paths.Paths {
		for method, operation := range path.Operations.Operations {
			// @TODO 删除
			//params := map[string]bool{
			//	//	//	"GetProjectAreaWeather_V1": true,
			//	//	//	"GetGoalSetting":           true,
			//	//	//	"ListAreaRanking":          true,
			//	//	//	"ListRole_V1":              true,
			//	//	//"ListActionAccess":               true,
			//	//	//"CreateForm":                     true,
			//	//	"CreateAction":     true,
			//"SaveStrategy_V1": true,
			//}
			//if params[operation.OperationId] {
			if operation.OperationId == "OpenAPI" || operation.OperationId == "ER" {
				continue
			}
			err := openApi.GenerateOperationDoc(nd, operation, url, method)
			if err != nil {
				return err
			}
			//}
		}
	}

	if err := openApi.doc.Validate(); err != nil {
		log.Fatalf("error during validation: %s", err)
		return err
	}
	openApi.doc.SaveToFile("iep-permission.docx")

	return nil
}

func (openApi *OpenApi) GenerateOperationDoc(nd document.NumberingDefinition, operation *oas.Operation, url string, method oas.HttpMethod) error {
	para := openApi.doc.AddParagraph()
	//para.SetNumberingDefinition(nd)
	para.Properties().SetHeadingLevel(4)
	para.AddRun().AddText(operation.Summary)

	openApi.doc.AddParagraph().AddRun().AddText("方法：")
	openApi.doc.AddParagraph().AddRun().AddText(fmt.Sprintf("    http://ip:port%s", url))
	openApi.doc.AddParagraph().AddRun().AddText("请求方式：")
	openApi.doc.AddParagraph().AddRun().AddText(fmt.Sprintf("    %s", Method(method)))
	if len(operation.Parameters) > 0 || operation.RequestBody != nil {
		err := openApi.GenerateOperationInput(operation)
		if err != nil {
			return err
		}
	}

	for statusCode, response := range operation.Responses.Responses {
		switch statusCode {
		case 200, 201:
			err := openApi.GenerateOperationOutput(response)
			if err != nil {
				return err
			}
		}
	}

	openApi.doc.AddParagraph()
	return nil
}

func (openApi *OpenApi) GenerateOperationOutput(response *oas.Response) error {
	openApi.doc.AddParagraph().AddRun().AddText("输出：")
	openApi.doc.AddParagraph().AddRun().AddText("    输出 JSON 格式的数据，格式如下")

	for contentType, resp := range response.ResponseObject.WithContent.Content {
		schema := resp.Schema
		switch contentType {
		case "application/json":
			err := openApi.GenerateOperationOutputStruct_JSON(schema)
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func (openApi *OpenApi) GenerateOperationOutputStruct_JSON(schema *oas.Schema) error {
	table := openApi.doc.AddTable()
	table.Properties().SetAlignment(wml.ST_JcTableCenter)
	table.Properties().SetWidthPercent(80)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)
	hdrRow := table.AddRow()

	cell := hdrRow.AddCell()
	cell.Properties().SetWidth(100 * measurement.Pixel72)
	cellPara := cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellPara.AddRun().AddText("参数名")

	cell = hdrRow.AddCell()
	cell.Properties().SetWidth(50 * measurement.Pixel72)

	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellPara.AddRun().AddText("返回类型")

	cell = hdrRow.AddCell()
	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellPara.AddRun().AddText("说明")

	_, err := openApi.GenerateOperationOutputTableAddRow(table, schema)
	if err != nil {
		return err
	}

	return nil
}

func (openApi *OpenApi) GenerateOperationOuthOutputSchemasProperties(referIDs ...string) error {
	isOut := map[string]bool{}
	for _, referID := range referIDs {
		if !isOut[referID] {
			isOut[referID] = true

			schema := openApi.api.Schemas[referID]

			if (len(schema.AllOf)+len(schema.Properties)) == 0 && schema.Refer == nil {
				continue
			}

			openApi.doc.AddParagraph().AddRun().AddText(fmt.Sprintf("%s:", referID))

			table := openApi.doc.AddTable()
			table.Properties().SetAlignment(wml.ST_JcTableCenter)
			table.Properties().SetWidthPercent(80)
			borders := table.Properties().Borders()
			borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)
			hdrRow := table.AddRow()

			cell := hdrRow.AddCell()
			cell.Properties().SetWidth(100 * measurement.Pixel72)
			cellPara := cell.AddParagraph()
			cellPara.Properties().SetAlignment(wml.ST_JcCenter)
			cellPara.AddRun().AddText("参数名")

			cell = hdrRow.AddCell()
			cellPara = cell.AddParagraph()
			cellPara.Properties().SetAlignment(wml.ST_JcCenter)
			cell.Properties().SetWidth(50 * measurement.Pixel72)
			cellPara.AddRun().AddText("返回类型")

			cell = hdrRow.AddCell()
			cellPara = cell.AddParagraph()
			cellPara.Properties().SetAlignment(wml.ST_JcCenter)
			cellPara.AddRun().AddText("说明")

			_, err := openApi.GenerateOperationOutputTableAddRow(table, schema)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (openApi *OpenApi) GenerateOperationOutputTableAddRow(table document.Table, schema *oas.Schema) ([]string, error) {
	referIDList := []string{}
	if _, ok := schema.Reference.Refer.(*oas.ComponentRefer); ok {
		refID := schema.Reference.Refer.(*oas.ComponentRefer).ID
		return openApi.GenerateOperationOutputTableAddRow(table, openApi.api.Schemas[refID])
	}

	if len(schema.Properties) > 0 {
		for param, propertie := range schema.Properties {

			switch propertie.Type {
			case oas.TypeObject:
				for k, v := range propertie.Properties {
					row := table.AddRow()
					row.AddCell().AddParagraph().AddRun().AddText(k)
					if len(v.AllOf) > 0 {
						if _, ok := v.AllOf[0].Reference.Refer.(*oas.ComponentRefer); ok {
							opID := v.AllOf[0].Reference.Refer.(*oas.ComponentRefer).ID
							referIDList = append(referIDList, opID)
							row.AddCell().AddParagraph().AddRun().AddText(Type(opID))
						}
					} else {
						row.AddCell().AddParagraph().AddRun().AddText(Type(string(v.Type)))
					}
					row.AddCell().AddParagraph().AddRun().AddText(v.Description)
				}
				break
			case oas.TypeArray:

				row := table.AddRow()
				row.AddCell().AddParagraph().AddRun().AddText(param)

				if _, ok := propertie.Items.Refer.(*oas.ComponentRefer); ok {
					opID := propertie.Items.Refer.(*oas.ComponentRefer).ID
					referIDList = append(referIDList, opID)
					row.AddCell().AddParagraph().AddRun().AddText("[]" + opID)
				} else {
					row.AddCell().AddParagraph().AddRun().AddText(Type(string(propertie.Type)))
				}
				row.AddCell().AddParagraph().AddRun().AddText(propertie.Description)

				break
			default:

				row := table.AddRow()
				row.AddCell().AddParagraph().AddRun().AddText(param)
				if len(propertie.AllOf) > 0 {
					if _, ok := propertie.AllOf[0].Reference.Refer.(*oas.ComponentRefer); ok {
						rID := propertie.AllOf[0].Reference.Refer.(*oas.ComponentRefer).ID
						row.AddCell().AddParagraph().AddRun().AddText(Type(rID))
						referIDList = append(referIDList, rID)
						row.AddCell().AddParagraph().AddRun().AddText(propertie.AllOf[1].Description)
					}
				} else {
					row.AddCell().AddParagraph().AddRun().AddText(Type(string(propertie.Type)))
					row.AddCell().AddParagraph().AddRun().AddText(propertie.Description)
				}
				break
			}

		}
		err := openApi.GenerateOperationOuthOutputSchemasProperties(referIDList...)
		if err != nil {
			return []string{}, err
		}
	}

	if len(schema.AllOf) > 0 {
		for _, proSchema := range schema.AllOf {
			ids, err := openApi.GenerateOperationOutputTableAddRow(table, proSchema)
			if err != nil {
				return ids, err
			}

			referIDList = append(referIDList, ids...)
		}
	}

	return referIDList, nil
}

func (openApi *OpenApi) GenerateOperationInput(operation *oas.Operation) error {
	openApi.doc.AddParagraph().AddRun().AddText("输入：")
	table := openApi.doc.AddTable()
	table.Properties().SetAlignment(wml.ST_JcTableCenter)
	table.Properties().SetWidthPercent(80)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)
	hdrRow := table.AddRow()

	cell := hdrRow.AddCell()
	cellPara := cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellPara.AddRun().AddText("参数名")

	cell = hdrRow.AddCell()
	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellPara.AddRun().AddText("是否必选")

	cell = hdrRow.AddCell()
	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellPara.AddRun().AddText("参数类型")

	cell = hdrRow.AddCell()
	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellPara.AddRun().AddText("说明")
	for _, param := range operation.Parameters {
		row := table.AddRow()
		row.AddCell().AddParagraph().AddRun().AddText(param.Name)
		if param.Required {
			row.AddCell().AddParagraph().AddRun().AddText("是")
		} else {
			row.AddCell().AddParagraph().AddRun().AddText("否")
		}
		row.AddCell().AddParagraph().AddRun().AddText(string(param.In))

		if len(param.Schema.AllOf) > 0 {
			if _, ok := param.Schema.AllOf[0].Reference.Refer.(*oas.ComponentRefer); ok {
				row.AddCell().AddParagraph().AddRun().AddText(param.Schema.AllOf[1].Description)
			} else {
				row.AddCell().AddParagraph().AddRun().AddText(param.Schema.Description)
			}
		} else {
			row.AddCell().AddParagraph().AddRun().AddText(param.Schema.Description)
		}
	}

	if operation.RequestBody != nil {
		for contentType, mediaType := range operation.RequestBody.Content {
			switch contentType {
			case "application/json":
				err := openApi.GenerateOperationInputRequestBody_JSON(mediaType)
				if err != nil {
					return err
				}
			}
		}

	}

	return nil
}

func (openApi *OpenApi) GenerateOperationInputRequestBody_JSON(mediaType *oas.MediaType) error {
	openApi.doc.AddParagraph().AddRun().AddText("Body JSON 输入：")

	jms := mediaType.MediaTypeObject.Schema
	if jms != nil {
		if len(jms.AllOf) > 0 {
			if _, ok := jms.AllOf[0].Reference.Refer.(*oas.ComponentRefer); ok {
				rID := jms.AllOf[0].Reference.Refer.(*oas.ComponentRefer).ID
				err := openApi.GenerateOperationInputBodyStruct(rID, nil, true)
				if err != nil {
					return err
				}

			}
		}
		if len(jms.Properties) > 0 {
			err := openApi.GenerateOperationInputBodyStruct("", jms, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (openApi *OpenApi) GenerateOperationInputBodyStruct(referID string, schema *oas.Schema, first bool, addStructList ...string) error {
	if referID != "" {
		schema = openApi.api.Schemas[referID]
	}

	if (len(schema.AllOf)+len(schema.Properties)) == 0 && schema.Refer == nil {
		return nil
	}
	if referID != "" || !first {
		openApi.doc.AddParagraph().AddRun().AddText(fmt.Sprintf("%s:", referID))
	}
	table := openApi.doc.AddTable()
	table.Properties().SetAlignment(wml.ST_JcTableCenter)
	table.Properties().SetWidthPercent(80)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)
	hdrRow := table.AddRow()

	cell := hdrRow.AddCell()
	cellPara := cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellPara.AddRun().AddText("参数名")

	cell = hdrRow.AddCell()
	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellPara.AddRun().AddText("是否必选")

	cell = hdrRow.AddCell()
	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellPara.AddRun().AddText("类型")

	cell = hdrRow.AddCell()
	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellPara.AddRun().AddText("说明")

	referIDList := []string{}
	if len(schema.Properties) > 0 {
		for key, propertie := range schema.Properties {
			row := table.AddRow()
			row.AddCell().AddParagraph().AddRun().AddText(key)
			isRequired := false
			for _, v := range schema.Required {
				if v == key {
					isRequired = true
				}
			}
			if isRequired {
				row.AddCell().AddParagraph().AddRun().AddText("是")
			} else {
				row.AddCell().AddParagraph().AddRun().AddText("否")
			}

			switch propertie.Type {
			case oas.TypeArray:
				if _, ok := propertie.Items.Refer.(*oas.ComponentRefer); ok {
					opID := propertie.Items.Refer.(*oas.ComponentRefer).ID
					referIDList = append(referIDList, opID)
					row.AddCell().AddParagraph().AddRun().AddText("[]" + opID)
				} else {
					row.AddCell().AddParagraph().AddRun().AddText("[]" + string(propertie.Items.Type))
				}

				break

			default:
				if len(propertie.AllOf) > 0 {
					opID := propertie.AllOf[0].Reference.Refer.(*oas.ComponentRefer).ID
					referIDList = append(referIDList, opID)
					row.AddCell().AddParagraph().AddRun().AddText(Type(opID))
				} else {
					row.AddCell().AddParagraph().AddRun().AddText(Type(string(propertie.Type)))
				}
			}

			row.AddCell().AddParagraph().AddRun().AddText(propertie.Description)
		}
	}

	if len(schema.AllOf) > 0 {
		err := openApi.GenerateOperationInputTableAddRow(table, schema.AllOf[0])
		if err != nil {
			return err
		}
	}

	err := openApi.GenerateOperationInputOuthSchemaTable(referIDList...)
	if err != nil {
		return err
	}

	return nil
}

func (openApi *OpenApi) GenerateOperationInputOuthSchemaTable(referIDs ...string) error {
	for _, referID := range referIDs {
		err := openApi.GenerateOperationInputBodyStruct(referID, nil, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (openApi *OpenApi) GenerateOperationInputTableAddRow(table document.Table, schema *oas.Schema) error {
	if _, ok := schema.Reference.Refer.(*oas.ComponentRefer); ok {
		refID := schema.Reference.Refer.(*oas.ComponentRefer).ID
		return openApi.GenerateOperationInputTableAddRow(table, openApi.api.Schemas[refID])
	}
	if len(schema.Properties) > 0 {
		for key, propertie := range schema.Properties {

			row := table.AddRow()
			row.AddCell().AddParagraph().AddRun().AddText(key)

			isRequired := false
			for _, v := range schema.Required {
				if v == key {
					isRequired = true
				}
			}
			if isRequired {
				row.AddCell().AddParagraph().AddRun().AddText("是")
			} else {
				row.AddCell().AddParagraph().AddRun().AddText("否")
			}
			row.AddCell().AddParagraph().AddRun().AddText(Type(string(propertie.Type)))
			if len(propertie.AllOf) > 0 {
				if _, ok := propertie.AllOf[0].Reference.Refer.(*oas.ComponentRefer); ok {
					row.AddCell().AddParagraph().AddRun().AddText(propertie.AllOf[1].Description)
				} else {
					row.AddCell().AddParagraph().AddRun().AddText(propertie.Description)
				}
			} else {
				row.AddCell().AddParagraph().AddRun().AddText(propertie.Description)
			}
		}
	}

	return nil
}

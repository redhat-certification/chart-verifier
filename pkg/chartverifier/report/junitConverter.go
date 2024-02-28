package report

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"time"
)

func encodeTokenArray(e *xml.Encoder, tokens []xml.Token) error {
	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}
	return e.Flush()
}

func encodeString(e *xml.Encoder, name string, value string) error {
	start := xml.StartElement{Name: xml.Name{"", name}}
	tokens := []xml.Token{start}
	tokens = append(tokens, xml.CharData(value), xml.EndElement{start.Name})

	return encodeTokenArray(e, tokens)
}

func encodeStringArray(e *xml.Encoder, name string, values []string) error {
	start := xml.StartElement{Name: xml.Name{"", name}}
	tokens := []xml.Token{start}

	for _, value := range values {
		t := xml.StartElement{Name: xml.Name{"", "value"}}
		tokens = append(tokens, t, xml.CharData(value), xml.EndElement{t.Name})
	}

	tokens = append(tokens, xml.EndElement{start.Name})

	return encodeTokenArray(e, tokens)
}

func encodeStringMap(e *xml.Encoder, name string, values map[string]string) error {
	start := xml.StartElement{Name: xml.Name{"", name}}
	tokens := []xml.Token{start}

	for key, value := range values {
		t := xml.StartElement{Name: xml.Name{"", key}}
		tokens = append(tokens, t, xml.CharData(value), xml.EndElement{t.Name})
	}

	tokens = append(tokens, xml.EndElement{start.Name})

	return encodeTokenArray(e, tokens)
}

func encodeReportMetadata(e *xml.Encoder, m ReportMetadata, start xml.StartElement) {
	// Start ReportMetadata and encode what can be done automatically
	e.EncodeToken(start)
	e.Encode(m.ToolMetadata)
	encodeString(e, "Overrides", m.Overrides)

	// Work through parts of the ChartData which is a problem due to Annotations field
	chartData := m.ChartData
	chartDataStartToken := xml.StartElement{Name: xml.Name{"", "ChartData"}}
	e.EncodeToken(chartDataStartToken)

	// Loop through helmchart.Metadata Fields and encode strings/bools
	v := reflect.ValueOf(*chartData)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldType := fmt.Sprintf("%T", v.Field(i).Interface())

		if fieldType == "string" || fieldType == "bool" {
			encodeString(e,
				typeOfS.Field(i).Name,
				fmt.Sprintf("%v", v.Field(i).Interface()))
		}
	}

	// Loop through helmchart.Metadata Fields
	encodeStringArray(e, "Sources", chartData.Sources)
	encodeStringArray(e, "Keywords", chartData.Keywords)
	encodeStringMap(e, "Annotations", chartData.Annotations)

	// Maintainers Start
	maintainersStartToken := xml.StartElement{Name: xml.Name{"", "Maintainers"}}
	e.EncodeToken(maintainersStartToken)
	e.Encode(chartData.Maintainers)
	e.EncodeToken(xml.EndElement{maintainersStartToken.Name})
	// Maintainers End

	// Dependencies Start
	dependenciesStartToken := xml.StartElement{Name: xml.Name{"", "Dependencies"}}
	e.EncodeToken(dependenciesStartToken)
	e.Encode(chartData.Dependencies)
	e.EncodeToken(xml.EndElement{dependenciesStartToken.Name})
	// Dependencies End

	e.EncodeToken(xml.EndElement{chartDataStartToken.Name})
	// ChartData End

	e.EncodeToken(xml.EndElement{start.Name})
	// Metadata End

	e.Flush()
}

func encodeResults(e *xml.Encoder, r Report) error {
	tokens := []xml.Token{}

	for _, testcase := range r.Results {

		// Put Check, Type, and Outcome into XML testcase element
		t := xml.StartElement{Name: xml.Name{"", "testcase"},
			Attr: []xml.Attr{
				{xml.Name{"", "name"}, string(testcase.Check)},
				{xml.Name{"", "classname"}, string(testcase.Reason)},
				{xml.Name{"", "assertions"}, "1"},
			},
		}

		//Create an element for Reason
		reasonToken := xml.StartElement{Name: xml.Name{"", "system-out"}}

		tokens = append(tokens, t,
			reasonToken, xml.CharData(testcase.Reason), xml.EndElement{reasonToken.Name},
			xml.EndElement{t.Name})

	}

	return encodeTokenArray(e, tokens)
}

func (r Report) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	numTests := len(r.Results)
	numFailures := 0
	numSkipped := 0
	numPassed := 0
	timestamp := time.Now().Format(time.RFC3339)

	for _, element := range r.Results {
		switch element.Outcome {
		case "FAIL":
			numFailures += 1
		case "SKIPPED":
			numSkipped += 1
		case "PASS":
			numPassed += 1
		}
	}

	junitStart := xml.StartElement{Name: xml.Name{"", "testsuites"},
		Attr: []xml.Attr{
			{xml.Name{"", "name"}, "chart-verifier test run"},
			{xml.Name{"", "tests"}, fmt.Sprint(numTests)},
			{xml.Name{"", "failures"}, fmt.Sprint(numFailures)},
			{xml.Name{"", "skipped"}, fmt.Sprint(numSkipped)},
			{xml.Name{"", "timestamp"}, timestamp},
		},
	}
	propertiesStart := xml.StartElement{Name: xml.Name{"", "properties"}}
	propertyStart := xml.StartElement{Name: xml.Name{"", "property"},
		Attr: []xml.Attr{
			{xml.Name{"", "name"}, "config"},
		},
	}

	e.EncodeToken(junitStart)
	e.EncodeToken(propertiesStart)
	e.EncodeToken(propertyStart)
	encodeString(e, "ApiVersion", r.Apiversion)
	encodeString(e, "Kind", r.Kind)
	e.Encode(r.options)
	encodeReportMetadata(e, r.Metadata, xml.StartElement{Name: xml.Name{"", "Metadata"}})
	e.EncodeToken(xml.EndElement{propertyStart.Name})
	e.EncodeToken(xml.EndElement{propertiesStart.Name})

	testSuiteStart := xml.StartElement{Name: xml.Name{"", "testsuite"},
		Attr: []xml.Attr{
			{xml.Name{"", "name"}, "chart-verifier test run"},
			{xml.Name{"", "tests"}, fmt.Sprint(numTests)},
			{xml.Name{"", "failures"}, fmt.Sprint(numFailures)},
			{xml.Name{"", "skipped"}, fmt.Sprint(numSkipped)},
			{xml.Name{"", "timestamp"}, timestamp},
		},
	}
	e.EncodeToken(testSuiteStart)

	// e.Encode(r.Results)
	encodeResults(e, r)

	// flush to ensure tokens are written
	e.EncodeToken(xml.EndElement{testSuiteStart.Name})
	e.EncodeToken(xml.EndElement{junitStart.Name})

	return e.Flush()
}

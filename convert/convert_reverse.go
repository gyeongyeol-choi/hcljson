package convert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	jsonParser "github.com/tmax-cloud/hcljson/parser"
	hclprinter "github.com/tmax-cloud/hcljson/printer"
	// hclprinter "github.com/hashicorp/hcl/hcl/printer"
	// jsonParser "github.com/hashicorp/hcl/json/parser"
)

func JsonToHcl(input []byte, typeSchemaStr string) []byte {

	var typeSchema map[string]interface{}
	json.Unmarshal([]byte(typeSchemaStr), &typeSchema)

	// MEMO: json 재구성 함수 호출
	input = regenJson(input)

	bytes, err := convertJsonToHcl(input, typeSchema)
	if err != nil {
		fmt.Errorf("hclTojson() error. %s", err)
	}
	return bytes
}

func convertJsonToHcl(input []byte, typeSchema map[string]interface{}) ([]byte, error) {
	ast, err := jsonParser.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse JSON: %s", err)
	}
	var buf bytes.Buffer
	if err := hclprinter.Fprint(&buf, ast, typeSchema); err != nil {
		return nil, fmt.Errorf("Unable to print HCL: %s", err)
	}

	return buf.Bytes(), nil
}

// MEMO : json 내 프로퍼티가 부모 프로퍼티 경로를 포함하도록 재구성 (map / object 구분 위함)
func regenJson(input []byte) []byte {

	var data map[string]interface{}
	json.Unmarshal(input, &data)

	for key, value := range data {
		if strings.Contains(key, "**##**") {
			continue
		}

		tmp := key
		key = tmp + "**##**"

		if reflect.TypeOf(value).Kind() == reflect.Map {
			data1 := value.(map[string]interface{})

			for key, value := range data1 {
				if strings.Contains(key, "**##**") {
					continue
				}

				tmp1 := key
				key = tmp + "**##**" + tmp1

				if reflect.TypeOf(value).Kind() == reflect.Map {
					data2 := value.(map[string]interface{})

					for key, value := range data2 {
						if strings.Contains(key, "**##**") {
							continue
						}

						tmp2 := key
						key = tmp + "**##**" + tmp1 + "**##**" + key

						if reflect.TypeOf(value).Kind() == reflect.Map {
							data3 := value.(map[string]interface{})

							for key, value := range data3 {
								if strings.Contains(key, "**##**") {
									continue
								}
								tmp3 := key
								key = tmp + "**##**" + tmp1 + "**##**" + tmp2 + "**##**" + key

								if reflect.TypeOf(value).Kind() == reflect.Map {
									data4 := value.(map[string]interface{})
									for key, value := range data4 {
										if strings.Contains(key, "**##**") {
											continue
										}
										tmp4 := key
										key = tmp + "**##**" + tmp1 + "**##**" + tmp2 + "**##**" + tmp3 + "**##**" + key

										if reflect.TypeOf(value).Kind() == reflect.Map {
											data5 := value.(map[string]interface{})
											for key, value := range data5 {
												if strings.Contains(key, "**##**") {
													continue
												}
												tmp5 := key
												key = tmp + "**##**" + tmp1 + "**##**" + tmp2 + "**##**" + tmp3 + "**##**" + tmp4 + "**##**" + key

												if reflect.TypeOf(value).Kind() == reflect.Map {
													data6 := value.(map[string]interface{})
													for key, value := range data6 {
														if strings.Contains(key, "**##**") {
															continue
														}

														tmp6 := key
														key = tmp + "**##**" + tmp1 + "**##**" + tmp2 + "**##**" + tmp3 + "**##**" + tmp4 + "**##**" + tmp5 + "**##**" + key

														data6[key] = value
														delete(data6, tmp6)

													}
													data5[key] = data6

												} else {
													data5[key] = value
												}
												delete(data5, tmp5)
											}
											data4[key] = data5

										} else {
											data4[key] = value
										}
										delete(data4, tmp4)
									}
									data3[key] = data4
								} else {
									data3[key] = value
								}
								delete(data3, tmp3)
							}
							//data2[key] = data3
							data2[tmp2] = data3
						} else {
							//data2[key] = value
							data2[tmp2] = value
						}
						//delete(data2, tmp2)
					}
					//data1[key] = data2
					data1[tmp1] = data2
				} else {
					//data1[key] = value
					data1[tmp1] = value
				}
				//delete(data1, tmp1)
			}
			//data[key] = data1
			data[tmp] = data1
		} else {
			//data[key] = value
			data[tmp] = value
		}
	}
	fmt.Println(data)
	output, _ := json.Marshal(data)
	return output
}

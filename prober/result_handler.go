package prober

import "dtool/utils"

func OutputHandler(data interface{}) (string, error) {
	var output_str string
	var err error
	switch value := data.(type) {
	case CacheStruct:
		result := make(map[string]map[int][]string)
		temp := make(map[int][]string)
		if len(value.dict) > 0 {
			for cache_id := range value.dict {
				for rdns := range value.dict[cache_id] {
					temp[cache_id] = append(temp[cache_id], rdns)
				}
			}
		}
		result[value.target] = temp
		output_str, err = utils.ToJSON(result, "")
	case RecursiveStruct:
		result := make(map[string][]string)
		temp := []string{}
		for rdns := range value.dict {
			temp = append(temp, rdns)
		}
		result[value.target] = temp
		output_str, err = utils.ToJSON(result, "")
	case SVCBResult:
		output_str, err = utils.ToJSON(data, "")
	}
	return output_str, err
}

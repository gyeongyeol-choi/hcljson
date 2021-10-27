package convert

import (
	"strings"
)

func HclToHcl(origin string, convert string) string {

	c_slice := strings.Split(convert, "\n")

	flag := false
	object := ""
	key := ""
	final_str := ""

	var m map[string]string
	m = make(map[string]string)

	for _, tmp := range c_slice {

		if strings.HasPrefix(tmp, "provider") || strings.HasPrefix(tmp, "module") ||
			strings.HasPrefix(tmp, "variable") || strings.HasPrefix(tmp, "resource") ||
			strings.HasPrefix(tmp, "data") {
			flag = true

			key = strings.Split(tmp, "{")[0]
			key = strings.TrimRight(key, " ")
			key = strings.Replace(key, "\"", "", -1)
			key = strings.Replace(key, " ", "_", -1)
		}

		if flag {
			object += tmp + "\n"

			if tmp == "}" {
				flag = false
				object = strings.Replace(object, "\n\n", "\n", -1)
				m[key] = object
				object = ""
			}
		}
	}

	flag = false
	object = ""

	var comment_map map[string]string
	comment_map = make(map[string]string)

	o_slice := strings.Split(origin, "\n")

	for _, tmp := range o_slice {
		tmp = strings.TrimSpace(tmp)

		if strings.HasPrefix(tmp, "provider") || strings.HasPrefix(tmp, "module") ||
			strings.HasPrefix(tmp, "variable") || strings.HasPrefix(tmp, "resource") ||
			strings.HasPrefix(tmp, "data") {
			flag = true

			key = strings.Split(tmp, "{")[0]
			key = strings.TrimRight(key, " ")
			key = strings.Replace(key, "\"", "", -1)
			key = strings.Replace(key, " ", "_", -1)
		}

		if !flag {
			final_str += tmp + "\n"
		} else {
			object += tmp + "\n"

			if strings.Contains(tmp, "#") {
				slice := strings.SplitN(tmp, "#", 2)
				cmt_val := "#" + slice[1]
				cmt_key := ""

				if strings.Contains(slice[0], "{") {
					cmt_key = "{"
				} else if strings.Contains(slice[0], "}") {
					cmt_key = "}"
				} else if strings.Contains(slice[0], "=") {
					cmt_key = strings.Split(slice[0], "=")[0]
					cmt_key = strings.TrimSpace(cmt_key)
				}
				comment_map[cmt_key] = cmt_val
			}

			if tmp == "}" {
				flag = false
				object = strings.Replace(object, "\n\n", "\n", -1)

				val, exists := m[key]

				if exists {
					slice := strings.Split(val, "\n")
					cmt_key := ""

					for _, tmp := range slice {
						if strings.Contains(tmp, "{") {
							cmt_key = "{"
						} else if strings.Contains(tmp, "}") {
							cmt_key = "}"
						} else if strings.Contains(tmp, "=") {
							cmt_key = strings.Split(tmp, "=")[0]
							cmt_key = strings.TrimSpace(cmt_key)
						}
						val2, exists2 := comment_map[cmt_key]

						if exists2 {
							final_str += tmp + val2 + "\n"
							cmt_key = ""
						} else {
							final_str += tmp + "\n"
						}
					}
				}
				object = ""
				key = ""
			}
		}
	}
	return final_str
}

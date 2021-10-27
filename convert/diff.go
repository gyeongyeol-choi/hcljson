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

			if strings.HasPrefix(tmp, "}") {
				flag = false
				object = strings.Replace(object, "\n\n", "\n", -1)

				m[key] = object
				object = ""
			}
		}

	}

	flag = false
	object = ""
	cnt := 0
	tmp_cmt := ""

	var comment_map map[string]string
	comment_map = make(map[string]string)

	var post_map map[string]string
	post_map = make(map[string]string)

	o_slice := strings.Split(origin, "\n")

	for _, tmp := range o_slice {
		tmp = strings.TrimSpace(tmp)

		if (strings.HasPrefix(tmp, "provider") || strings.HasPrefix(tmp, "module") ||
			strings.HasPrefix(tmp, "variable") || strings.HasPrefix(tmp, "resource") ||
			strings.HasPrefix(tmp, "data")) && !flag {
			flag = true

			key = strings.Split(tmp, "{")[0]
			key = strings.TrimRight(key, " ")
			key = strings.Replace(key, "\"", "", -1)
			key = strings.Replace(key, " ", "_", -1)

		}

		if !flag {
			final_str += tmp + "\n"
		} else {
			if strings.Contains(tmp, "{") {
				cnt++
			}
			if strings.Contains(tmp, "}") {
				cnt--
			}

			object += tmp + "\n"

			if strings.Contains(tmp, "=") {
				post_key := strings.Split(tmp, "=")[0]
				post_key = strings.TrimSpace(post_key)

				if !strings.Contains(post_key, "#") && tmp_cmt != "" {
					post_map[post_key] = tmp_cmt
					tmp_cmt = ""
				}
			}

			if strings.HasPrefix(tmp, "#") && cnt == 1 {
				tmp_cmt += tmp + "\n"
			}

			if strings.Contains(tmp, "#") && cnt == 1 {
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

				if cmt_key != "" {
					comment_map[cmt_key] = cmt_val
				}
			}

			if strings.HasPrefix(tmp, "}") && cnt == 0 {

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

						val2, exists2 := post_map[cmt_key]

						if exists2 {
							spaceCnt := countLeadingSpaces(tmp)

							for i := 0; i < spaceCnt; i++ {
								final_str += " "
							}
							final_str += val2
						}
						val3, exists3 := comment_map[cmt_key]

						if exists3 {
							final_str += tmp + val3 + "\n"
						} else {
							final_str += tmp + "\n"
						}

						cmt_key = ""
					}
				}
				object = ""
				key = ""
			}
		}
	}
	return final_str
}

func countLeadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}
